package controller

import (
	"blue-bell_back/models"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

// InitTrans 初始化翻译器
// 参数: locale - 语言环境
// 返回值: 错误对象，如果初始化失败则返回错误
func InitTrans(locale string) (err error) {
	//修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		//注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

			if name == "-" {
				return ""
			}
			return name
		})

		//为SignUpPrarm注册自定义校验方法
		v.RegisterStructValidation(SignUpParamStructLevelValidation, models.ParamSignUp{})

		zhT := zh.New() //中文
		enT := en.New() //英文翻译器

		//初始化统一翻译器
		uni := ut.New(enT, zhT, enT)

		//获取指定语言环境的翻译器
		var ok bool
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		//注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
	}
	return
}

// removeTopStruct 去除提示信息中的结构体名称
// 参数: fields - 包含错误提示信息的映射
// 返回值: 去除结构体名称后的错误提示信息映射

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

// SignUpParamStructLevelValidation 自定义SignUpParam结构体校验函数
// 参数: sl - 结构体级别的验证器
func SignUpParamStructLevelValidation(sl validator.StructLevel) {
	su := sl.Current().Interface().(models.ParamSignUp)

	//检验密码是否和实际的一致
	if su.Password != su.RePassword {
		// 输出错误提示信息，最后一个参数就是传递的param
		sl.ReportError(su.RePassword, "re_password",
			"RePassword", "eqfield", "password")
	}

}
