package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/guonaihong/gout"
)

var QYWechatGroupBotWebHookURL string = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=70f5d200-03dc-4a29-9a72-f60a1d7ec9e6"

func main() {
	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	result := ""
	for j := 0; j < len(output); j++ {
		result += string(output[j])
	}

	if strings.Contains(result, "失败用例：") {
		sendFailMessageToQYWX(result)
	}
}

//发送到企业微信群
func sendFailMessageToQYWX(result string) {
	lines := strings.Split(result, "\n")

	t := time.Now()
	output := t.Format("[01月02日 15:04:05有故障！] \n")
	for _, tmpItem := range lines {
		isCaseIndex, err := regexp.MatchString(`\(\d+/\d+\)`, tmpItem)
		if err != nil {
			fmt.Print(err)
		}

		if isCaseIndex {
			if strings.Contains(tmpItem, "失败") {
				//(1/1) 失败 [xxx.sh] [0. xx] (2.00s)
				regexp, err := regexp.Compile(`\(\d+/\d+\) 失败 \[.*\]\s\[\d+\.\s(.*)\].*`)
				if err != nil {
					fmt.Println(err)
				}
				match := regexp.FindStringSubmatch(tmpItem)

				output += `> <font color="warning">` + match[1] + "</font> [x] \n"
			}
		}
	}

	err := gout.POST(QYWechatGroupBotWebHookURL).
		SetJSON(gout.H{
			"msgtype": "markdown",
			"markdown": gout.H{
				"content": output,
			}}).Do()

	if err != nil {
		fmt.Println(err)
	}
}
