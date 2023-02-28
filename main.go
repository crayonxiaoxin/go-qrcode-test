package main

import (
	"bytes"
	"flag"
	"fmt"
	"go-qrcode/utils"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func main() {
	r := gin.Default()
	r.GET("/", qrcHandler)
	r.POST("/", qrcHandler)

	// 监听端口（如需更改端口，运行 go run . -port 8086）
	port := flag.Int("port", 8085, "端口")
	// 解析参数
	flag.Parse()
	// 运行
	r.Run(":" + strconv.Itoa(*port))
}

// 示例
//
// http://localhost:8085/?content=https://example.com&bg=&fg=&block=&border=
func qrcHandler(ctx *gin.Context) {

	// 内容
	content := ctx.Query("content")
	if content == "" {
		ctx.JSON(500, gin.H{"error": "Content could not be empty."})
		return
	}

	// size
	blockWidthStr := ctx.Query("block")
	borderWidthStr := ctx.Query("border")
	blockWidth := validSize(blockWidthStr, 10)
	borderWidth := validSize(borderWidthStr, 20)

	// 颜色
	bgColorStr := ctx.Query("bg")
	fgColorStr := ctx.Query("fg")
	backgroundColorHex := validHexColor(bgColorStr, "#ffffff")
	foregroundColorHex := validHexColor(fgColorStr, "#000000")

	// Logo
	logo := ""                       // 保存位置
	ext := ""                        // 图片格式
	dstDir := "uploads"              // 保存目录
	fh, err3 := ctx.FormFile("logo") // 上传文件
	if err3 == nil {
		if !utils.PathExists(dstDir) {
			err4 := os.Mkdir(dstDir, os.ModePerm)
			if err4 != nil {
				ctx.JSON(500, gin.H{"error": err4})
				return
			}
		}
		// 文件名
		filename := filepath.Base(fh.Filename)
		// 后缀
		ext = filepath.Ext(filename)
		if isJPG(ext) || isPNG(ext) {
			// 重命名
			newName := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
			dst := fmt.Sprintf("%v/%v", dstDir, newName)
			err4 := utils.SaveFile(fh, dst)
			if err4 == nil {
				logo = dst
			} else {
				ctx.JSON(500, gin.H{"error": err4})
				return
			}
		}
	}

	qr, err := qrcode.NewWith(content, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err})
		return
	}

	// 写入器
	buf := bytes.NewBuffer(nil)
	wr := nopCloser{Writer: buf}
	var w *standard.Writer
	hasLogo := logo != "" && ext != ""
	// 生成二维码
	w = standard.NewWithWriter(
		wr,
		standard.WithQRWidth(uint8(blockWidth)),
		standard.WithBorderWidth(borderWidth),
		standard.WithBgColorRGBHex(backgroundColorHex),
		standard.WithFgColorRGBHex(foregroundColorHex),
	)
	if hasLogo { // 有 logo
		// 获取二维码大小
		qrWidth := w.Attribute(qr.Dimension()).W
		// 调整 logo 大小
		logoWidth := qrWidth / 5
		i, err4 := utils.OpenImg(logo)
		if err4 != nil {
			ctx.JSON(500, gin.H{"error": err4})
			return
		} else {
			i2, err5 := utils.ResizeImg(i, logoWidth)
			if err5 != nil {
				ctx.JSON(500, gin.H{"error": err5})
				return
			} else {
				w = standard.NewWithWriter(
					wr,
					standard.WithQRWidth(uint8(blockWidth)),
					standard.WithBorderWidth(borderWidth),
					standard.WithBgColorRGBHex(backgroundColorHex),
					standard.WithFgColorRGBHex(foregroundColorHex),
					standard.WithLogoImage(i2),
				)
			}
		}
	}

	// 保存
	err2 := qr.Save(w)
	if err2 != nil {
		ctx.JSON(500, gin.H{"error": err2})
		return
	}

	// 移除logo文件
	if hasLogo {
		utils.RemoveFile(logo)
	}

	// 输出到浏览器
	png := buf.Bytes()
	ctx.Writer.Header().Set("Content-Type", "image/png")
	ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(png)))
	ctx.Writer.Write(png)
}

// 纠正 hex 颜色格式
func validHexColor(bgColorStr string, defaultVal ...string) string {
	bgColorLen := len(bgColorStr)
	if bgColorLen >= 3 && bgColorLen <= 7 {
		if !strings.Contains(bgColorStr, "#") {
			bgColorStr = "#" + bgColorStr
		}
		if len(bgColorStr) == 4 || len(bgColorStr) == 7 {
			return bgColorStr
		}
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return ""
}

func isPNG(ext string) bool {
	return ext == ".png"
}

func isJPG(ext string) bool {
	return ext == ".jpeg" || ext == ".jpg"
}

// 纠正 int 参数格式
func validSize(sizeStr string, defaultVal ...int) int {
	i, err := strconv.Atoi(sizeStr)
	if err == nil {
		return i
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return 0
}

// 写入器
type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error {
	return nil
}
