package main

import (
	"fmt"

	//lint:ignore ST1001 我需要 dot imports 来简化调用
	. "sing-box-sub-helper/packages"
)

func main() {
	fmt.Println("Sing-Box Subscription Helper")
	fmt.Println("Version: 0.3.0")
	fmt.Println("============================================================")
	fmt.Println("获取配置信息...")
	settings, err := GetSettings()
	if err != nil {
		fmt.Println("获取配置失败:", err)
		return
	}
	// 如果settings.StartStep小于等于1，从远程获取开始
	if settings.StartStep <= 1 {
		fmt.Println("从远程获取订阅数据...")
		base64Data, err := FetchBase64Data(settings.SubscribeURL)
		if err != nil {
			fmt.Println("获取订阅数据失败:", err)
			return
		}
		fmt.Println("写入base64数据文件...")
		err = SeveFile(settings.Base64File, base64Data)
		if err != nil {
			fmt.Println("写入base64数据文件失败:", err)
			return
		}
	}
	// 如果settings.StartStep小于等于2
	if settings.StartStep <= 2 {
		fmt.Println("读取base64数据文件并解密...")
		listData, err := ReadBase64FileDecode(settings.Base64File)
		if err != nil {
			fmt.Println("读取base64文件失败:", err)
		}
		fmt.Println("写入临时节点列表文件...")
		err = SeveFile(settings.TempListPath, listData)
		if err != nil {
			fmt.Println("写入temp.list失败:", err)
			return
		}
	}
	// 如果settings.StartStep小于等于3
	if settings.StartStep <= 3 {
		fmt.Println("转换订阅列表为json数据，并按过滤器过滤...")
		listData, err := ConvertSubscriptionToJson(settings.TempListPath, settings.Filter)
		if err != nil {
			fmt.Println("转换订阅列表失败:", err)
			return
		}
		fmt.Println("写入临时json文件...")
		err = SeveFile(settings.TempJsonPath, listData)
		if err != nil {
			fmt.Println("写入temp.json失败:", err)
			return
		}
	}

	// 合并到模板文件
	fmt.Println("合并模板文件...")
	err = MergeTemplateWithSubscription(settings.TemplatePath, settings.TempJsonPath, settings.OutputPath)
	if err != nil {
		fmt.Println("合并模板文件失败:", err)
		return
	}
	fmt.Println("成功运行结束，输出文件：", settings.OutputPath)
}
