package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"os"
)

func main() {
	var src, dst, to string
	var gov bool
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "src",
				Value:       "",
				Usage:       "yml, yaml file or json file to convert. ",
				Destination: &src,
			},
			&cli.StringFlag{
				Name:        "dst",
				Value:       "",
				Usage:       "convert to yml, yaml file or json file. ",
				Destination: &dst,
			},
			&cli.StringFlag{
				Name:        "to",
				Value:       "",
				Usage:       "convert dst type,can by yml or json.",
				Destination: &to,
			},
			&cli.BoolFlag{
				Name:        "var",
				Value:       true,
				Usage:       "convert dst to golang var.",
				Destination: &gov,
			},
		},
		Action: func(c *cli.Context) error {
			if to == "json" {
				toJson(src, dst, gov)
			} else if to == "yml" || to == "yaml" {
				toYaml(src, dst)
			}
			return nil
		},
	}

	app.Run(os.Args)
}

func toJson(src, dst string, gov bool) {
	buffer, err := ioutil.ReadFile(src)
	j2, err := yaml.YAMLToJSON(buffer)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	if gov {
		createVarFile(string(j2))
	}

	var mf *os.File
	if checkFileIsExist(dst) { //如果文件存在
		mf, err = os.OpenFile(dst, os.O_TRUNC|os.O_CREATE, os.ModePerm) //打开文件
	} else {
		mf, err = os.Create(dst) //创建文件
	}
	if err != nil {
		fmt.Printf("不能创建或者打开文件：%v\n", err)
		return
	}
	defer mf.Close()

	mf.WriteString(string(j2))

}

func toYaml(src, dst string) {
	buffer, err := ioutil.ReadFile(src)

	j2, err := yaml.JSONToYAML(buffer)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	var mf *os.File
	if checkFileIsExist(dst) { //如果文件存在
		mf, err = os.OpenFile(dst, os.O_TRUNC|os.O_CREATE, os.ModePerm) //打开文件
	} else {
		mf, err = os.Create(dst) //创建文件
	}
	if err != nil {
		fmt.Printf("不能创建或者打开文件：%v\n", err)
		return
	}
	defer mf.Close()

	mf.WriteString(string(j2))

}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func createVarFile(s string) {
	mfile := "docs/docsvar.go"
	var mf *os.File
	var err error
	if checkFileIsExist(mfile) { //如果文件存在
		mf, err = os.OpenFile(mfile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm) //打开文件
	} else {
		mf, err = os.Create(mfile) //创建文件
	}
	if err != nil {
		fmt.Printf("不能创建或者打开文件：%v\n", err)
		return
	}
	defer mf.Close()

	var c = " package docs"
	var c1 = "	var doc = ` "

	_, err = mf.WriteString(c)
	if err != nil {
		fmt.Printf("不能写入文件：%v\n", err)
		return
	}
	mf.WriteString("\n")
	mf.WriteString(c1)
	mf.WriteString(s)
	mf.WriteString("`")
}
