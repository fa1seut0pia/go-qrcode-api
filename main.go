package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"unicode"
)

func main() {
	http.HandleFunc("/qrcode", handleQRCodeRequest)
	log.Println("Starting server on :7688")
	if err := http.ListenAndServe(":7688", nil); err != nil {
		log.Fatal(err)
	}
}

func handleQRCodeRequest(w http.ResponseWriter, r *http.Request) {

	log.Println("Handling request for /qrcode?" + r.URL.RawQuery)

	text := r.URL.Query().Get("t")
	if text == "" {
		http.Error(w, "Parameter 't' is required", http.StatusBadRequest)
		return
	}

	base64Flag := r.URL.Query().Get("e")
	if base64Flag != "" {
		decodedText, err := base64.StdEncoding.DecodeString(text)
		if err != nil {
			http.Error(w, "Failed to decode base64 text", http.StatusBadRequest)
			return
		}
		text = string(decodedText)
	}

	qrCode, err := generateQRCode(text)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	logoURL := r.URL.Query().Get("l")
	if logoURL != "" {
		err := addLogoToQRCode(&qrCode, logoURL)
		if err != nil {
			http.Error(w, "Failed to add logo to QR code", http.StatusInternalServerError)
			return
		}
	}

	watermark := r.URL.Query().Get("w")
	if watermark != "" {
		err := addWatermark(&qrCode, watermark)
		if err != nil {
			http.Error(w, "Failed to add watermark", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "image/png")
	_ = png.Encode(w, qrCode)
}

func generateQRCode(text string) (image.Image, error) {
	qrCode, err := qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewReader(qrCode))
}

func addLogoToQRCode(qrCode *image.Image, logoURL string) error {
	resp, err := http.Get(logoURL)
	if err != nil {
		return fmt.Errorf("failed to fetch logo: %w", err)
	}
	defer resp.Body.Close()

	logo, _, err := image.Decode(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to decode logo: %w", err)
	}

	b := (*qrCode).Bounds()
	qrWidth := b.Dx()
	qrHeight := b.Dy()
	logoWidth := qrWidth / 6
	logoHeight := int(float64(logoWidth) * float64(logo.Bounds().Dy()) / float64(logo.Bounds().Dx()))
	resizedLogo := resize.Resize(uint(logoWidth), uint(logoHeight), logo, resize.Lanczos3)

	offset := image.Pt((qrWidth-logoWidth)/2, (qrHeight-logoHeight)/2)

	finalImage := image.NewRGBA(image.Rect(0, 0, qrWidth, qrHeight))
	draw.Draw(finalImage, finalImage.Bounds(), *qrCode, image.Point{}, draw.Over)
	draw.Draw(finalImage, resizedLogo.Bounds().Add(offset), resizedLogo, image.Point{}, draw.Over)

	*qrCode = finalImage

	return nil
}

func addWatermark(qrCode *image.Image, watermarkText string) error {

	fontBytes, err := os.ReadFile("SarasaFixedSC-Regular.ttf")
	if err != nil {
		return fmt.Errorf("failed to load font: %w", err)
	}

	trueTypeFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return fmt.Errorf("failed to parse font: %w", err)
	}

	b := (*qrCode).Bounds()

	const fontSize = 14.0
	const dpi = 72

	finalImage := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()+fontSize))
	draw.Draw(finalImage, finalImage.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(finalImage, b, *qrCode, image.Point{}, draw.Over)

	ctx := freetype.NewContext()
	ctx.SetFont(trueTypeFont)
	ctx.SetFontSize(fontSize)
	ctx.SetDst(finalImage)
	ctx.SetClip(finalImage.Bounds())
	ctx.SetDPI(dpi)

	chineseCount, englishCount := countChineseAndEnglish(watermarkText)
	textWidth := fontSize * (chineseCount + englishCount/2)
	if textWidth > b.Dx() {
		textWidth = b.Dx()
	}

	const textHeight = fontSize

	x := (b.Dx() - textWidth) / 2
	y := finalImage.Bounds().Dy() - textHeight

	watermarkColor := color.Black
	ctx.SetSrc(image.NewUniform(watermarkColor))
	_, err = ctx.DrawString(watermarkText, freetype.Pt(x, y))
	if err != nil {
		return fmt.Errorf("failed to draw watermark: %w", err)
	}

	*qrCode = finalImage

	return nil
}

func countChineseAndEnglish(s string) (int, int) {
	var chineseCount, englishCount int
	for _, r := range s {
		if unicode.Is(unicode.Scripts["Han"], r) {
			chineseCount++
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) && r < unicode.MaxASCII {
			englishCount++
		}
	}
	return chineseCount, englishCount
}
