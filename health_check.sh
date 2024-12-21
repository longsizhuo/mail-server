#!/bin/bash

# 配置部分
URL="http://localhost:8181/health"   # 要检查的健康检查 URL
CHECK_INTERVAL=600                    # 检查间隔（秒）
FAIL_THRESHOLD=5                   # 允许的最大连续失败次数
FAILED_COUNT=0                      # 初始化失败计数
EMAIL_TO="longsizhuo@gmail.com"      # 收件人邮箱
EMAIL_SUBJECT="服务健康检查告警"
EMAIL_BODY="警告：服务在过去 10 分钟内无法访问，请立即检查并重启！"

# 检查端口是否被占用的函数
close_port_8181() {
    echo "正在检查并关闭使用 8181 端口的进程..."
    # 查找占用 8181 端口的进程并强制杀掉
    PID=$(lsof -t -i:8181)
    if [ -n "$PID" ]; then
        sudo kill -9 "$PID"
        echo "端口 8181 已被释放 (PID: $PID)"
    else
        echo "端口 8181 未被占用"
    fi
}

# 重启 mail-server 服务的函数
restart_mail_server() {
    echo "正在尝试重启 mail-server 服务..."
    sudo systemctl restart mail-server
    if [ $? -eq 0 ]; then
        echo "mail-server 服务重启成功"
    else
        echo "mail-server 服务重启失败"
        return 1  # 返回 1 表示重启失败
    fi
}

while true; do
    # 发送请求并获取 HTTP 状态码
    HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$URL")

    if [ "$HTTP_STATUS" -eq 200 ]; then
        # 如果成功，重置失败计数
        FAILED_COUNT=0
        echo "$(date): 服务正常，状态码 $HTTP_STATUS"
    else
        # 如果失败，增加失败计数
        ((FAILED_COUNT++))
        echo "$(date): 服务异常，状态码 $HTTP_STATUS，连续失败次数：$FAILED_COUNT"
    fi

    # 检查失败计数是否达到阈值
    if [ "$FAILED_COUNT" -ge "$FAIL_THRESHOLD" ]; then
        echo "$(date): 检测到服务异常，尝试重启 mail-server"

        # 关闭 8181 端口进程
        close_port_8181

        # 尝试重启 mail-server 服务
        restart_mail_server

        # 重启后再次检查服务状态
        sleep 5  # 等待几秒钟，确保服务完全启动
        NEW_HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$URL")

        if [ "$NEW_HTTP_STATUS" -eq 200 ]; then
            # 服务恢复正常，重置失败计数，不发送邮件
            FAILED_COUNT=0
            echo "$(date): mail-server 已成功恢复，状态码 $NEW_HTTP_STATUS，不发送告警邮件。"
        else
            # 服务未恢复，发送告警邮件
            echo "$EMAIL_BODY" | mail -s "$EMAIL_SUBJECT" "$EMAIL_TO"
            echo "$(date): 服务仍然异常，已发送告警邮件至 $EMAIL_TO"
            # 重置失败计数，避免重复发送邮件
            FAILED_COUNT=0
        fi
    fi

    # 等待指定的检查间隔
    sleep "$CHECK_INTERVAL"
done
