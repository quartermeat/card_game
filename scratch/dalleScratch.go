package scratch

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/Andrew-peng/go-dalle2/dalle2"
	"github.com/joho/godotenv"
	"github.com/lxn/walk"
	decl "github.com/lxn/walk/declarative"
	"golang.org/x/sys/windows"
)

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func RunDalleTest() {

	var inputTextEdit *walk.TextEdit
	var runButton *walk.PushButton
	var openFileDialogButton *walk.PushButton
	var imagePath, maskPath string

	windowSize := decl.Size{
		Width:  400,
		Height: 300,
	}

	enableRunButton := func() {
		runButton.SetEnabled(len(inputTextEdit.Text()) > 0 && len(imagePath) > 0 && len(maskPath) > 0)
	}

	fmt.Println("hello world")

	var mainWindowPtr *walk.MainWindow

	mainWindow := &decl.MainWindow{
		AssignTo: &mainWindowPtr,
		Title:    "DALLE Image Editor",
		MinSize:  windowSize,
		Layout:   decl.VBox{},
		Children: []decl.Widget{
			decl.Label{
				Text: "Enter a new description:",
			},
			decl.TextEdit{
				AssignTo:      &inputTextEdit,
				OnTextChanged: enableRunButton,
			},
			decl.PushButton{
				AssignTo: &openFileDialogButton,
				Text:     "Select Image and Mask",
				OnClicked: func() {
					imageFilter := "Image Files (*.png, *.jpg, *.jpeg)|*.png;*.jpg;*.jpeg"

					dlg := new(walk.FileDialog)
					dlg.Title = "Select Image"
					dlg.Filter = imageFilter

					if ok, err := dlg.ShowOpen(nil); err != nil {
						log.Printf("Error opening image file: %v", err)
					} else if !ok {
						log.Println("Image selection canceled")
					} else {
						imagePath = dlg.FilePath
					}

					dlg.Title = "Select Mask"
					if ok, err := dlg.ShowOpen(nil); err != nil {
						log.Printf("Error opening mask file: %v", err)
					} else if !ok {
						log.Println("Mask selection canceled")
					} else {
						maskPath = dlg.FilePath
					}
				},
			},
			decl.PushButton{
				AssignTo: &runButton,
				Text:     "Run",
				Enabled:  false,
				OnClicked: func() {
					description := inputTextEdit.Text()
					RunDalle(description, imagePath, maskPath)
				},
			},
		},
	}

	if _, err := mainWindow.Run(); err != nil {
		log.Fatal(err)
	}

	user32 := windows.NewLazySystemDLL("user32.dll")
	setForegroundWindow := user32.NewProc("SetForegroundWindow")
	setForegroundWindow.Call(uintptr(mainWindowPtr.Handle()))

	fmt.Println("done running")
}

func RunDalle(description string, imagePath string, maskPath string) {
	// get secrets
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("DALLE_API_KEY")

	// connect a client to dalle 2
	client, err := dalle2.MakeNewClientV1(apiKey)
	if err != nil {
		log.Fatalf("Error initializing client: %s", err)
	}

	// read in an image
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		fmt.Printf("Error reading image file: %v\n", err)
		os.Exit(1)
	}
	maskData, err := os.ReadFile(maskPath)
	if err != nil {
		fmt.Printf("Error reading image file: %v\n", err)
		os.Exit(1)
	}

	resp, err := client.Edit(
		context.Background(),
		imageData,
		maskData,
		description,
		dalle2.WithNumImages(1),
		dalle2.WithSize(dalle2.LARGE),
		dalle2.WithFormat(dalle2.URL),
	)

	if err != nil {
		log.Fatal(err)
	}
	for _, img := range resp.Data {
		openbrowser(img.Url)
	}
}
