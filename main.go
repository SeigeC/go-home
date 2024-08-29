package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	// 获取今天的解锁时间，不包括毫秒和时区偏移量
	dateCmd := `log show --predicate='(process == "loginwindow") AND (composedMessage endswith "unlock success")' --info --style syslog --start "$(date '+%Y-%m-%d') 00:00:00" --end "$(date '+%Y-%m-%d') 23:59:59" | grep 'unlock success' | head -n 1 | awk '{print $1, substr($2,1,8)}'`

	cmd := exec.Command("sh", "-c", dateCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	// 提取解锁时间
	unlockTimeStr := strings.TrimSpace(out.String())
	fmt.Println("Unlock time:", unlockTimeStr)

	// 解析解锁时间
	unlockTime, err := time.Parse(time.DateTime, unlockTimeStr)
	if err != nil {
		fmt.Println("Error parsing unlock time:", err)
		return
	}

	// 计算 9 小时后的时间
	reminderTime := unlockTime.Add(9 * time.Hour)
	fmt.Println("提醒时间：", reminderTime.Format(time.DateTime))

	<-time.After(time.Until(reminderTime))

	// 设置提醒
	osascriptCmd := fmt.Sprintf(`osascript -e 'display notification "%s %s - %s 已工作满 8 小时" with title "准备下班"'`, unlockTime.Format(time.DateOnly), unlockTime.Format(time.Kitchen), reminderTime.Format(time.Kitchen))
	cmd = exec.Command("sh", "-c", osascriptCmd)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error setting reminder:", err)
		return
	}
}
