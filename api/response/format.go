package response

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 递归处理数据结构，转换时间和ID字段
func processData(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	value := reflect.ValueOf(data)

	// 处理指针
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		return processData(value.Elem().Interface())
	}

	switch value.Kind() {
	case reflect.Struct:
		// 处理time.Time类型
		if t, ok := data.(time.Time); ok {
			return t.Unix()
		}

		// 创建一个新的map来存储处理后的结构体字段
		result := make(map[string]interface{})

		// 处理结构体的字段
		t := value.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			// 跳过未导出的字段
			if !field.IsExported() {
				continue
			}

			// 获取json标签
			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}

			// 解析json标签获取字段名
			jsonName := field.Name
			if jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				if parts[0] != "" {
					jsonName = parts[0]
				}
			}

			fieldValue := value.Field(i).Interface()

			// 特别处理gorm.Model
			if field.Name == "Model" {
				// 检查是否为gorm.Model类型
				modelType := reflect.TypeOf(gorm.Model{})
				if reflect.TypeOf(fieldValue).ConvertibleTo(modelType) {
					model := fieldValue.(gorm.Model)
					// 添加小写id字段
					result["id"] = model.ID
					// 转换时间字段为时间戳
					result["created_at"] = model.CreatedAt.Unix()
					result["updated_at"] = model.UpdatedAt.Unix()
					if !model.DeletedAt.Time.IsZero() {
						result["deleted_at"] = model.DeletedAt.Time.Unix()
					}
					continue
				}
			}

			// 递归处理其他字段
			result[jsonName] = processData(fieldValue)
		}
		return result

	case reflect.Slice, reflect.Array:
		// 处理切片和数组
		resultSlice := make([]interface{}, value.Len())
		for i := 0; i < value.Len(); i++ {
			resultSlice[i] = processData(value.Index(i).Interface())
		}
		return resultSlice

	case reflect.Map:
		// 处理映射
		resultMap := make(map[string]interface{})
		keys := value.MapKeys()
		for _, key := range keys {
			keyStr := fmt.Sprintf("%v", key.Interface())
			resultMap[keyStr] = processData(value.MapIndex(key).Interface())
		}
		return resultMap

	default:
		// 其他类型直接返回
		return data
	}
}

func ReturnErrorWithData(c *gin.Context, data responseData, result interface{}) {
	data.Timestamp = time.Now().Unix()
	data.Data = processData(result)
	c.JSON(http.StatusOK, data)
	// Return directly
	c.Abort()
}

// ResponseData 正常响应
func ReturnData(c *gin.Context, result interface{}) {
	data := Success
	data.Timestamp = time.Now().Unix()
	data.Data = processData(result)
	c.JSON(http.StatusOK, data)
	// Return directly
	c.Abort()
}

// ResponseDataWithCount 正常响应
func ReturnDataWithCount(c *gin.Context, count int, result interface{}) {
	data := Success
	data.Timestamp = time.Now().Unix()
	data.Data = processData(result)
	data.Count = &count
	c.JSON(http.StatusOK, data)
	// Return directly
	c.Abort()
}

// ResponseError 错误响应
func ReturnError(c *gin.Context, data responseData, description string) {
	data.Timestamp = time.Now().Unix()
	data.Message = func() string {
		if description == "" {
			return data.Message
		}
		return description
	}()
	c.JSON(http.StatusOK, data)
	// Return directly
	c.Abort()
}

// ResponseSuccess 执行成功
func ReturnSuccess(c *gin.Context) {
	data := Success
	data.Timestamp = time.Now().Unix()
	c.JSON(http.StatusOK, data)
	// Return directly
	c.Abort()
}
