package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/quantstop/quantstopterminal/assets/images"
	"github.com/quantstop/quantstopterminal/internal/config"
	"github.com/quantstop/quantstopterminal/internal/engine"
	qstlog "github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/pkg/system"
	"github.com/skratchdot/open-golang/open"
	"log"
	"path/filepath"
	"strings"
)

var (
	BuildFlagVersion   string         // Build flag for version
	BuildFlagIsRelease string         // Build flag for setting release blurb
	Engine             *engine.Engine // Global pointer to Engine
	Config             *config.Config // Global pointer to Config

)

func main() {

	var err error

	// Create config
	Config = &config.Config{}

	// Setup config
	if err = Config.SetupConfig(); err != nil {
		log.Fatalf("Error settup up config: %s\n", err)
	}

	// Verify config
	if err = Config.CheckConfig(); err != nil {
		log.Fatalf("Error checking config: %s\n", err)
	}
	// Setup global logger
	if err = qstlog.SetupGlobalLogger(); err != nil {
		log.Fatalf("Error setting up global logger: %s\n", err)
	}

	// Setup all sub loggers
	if err = qstlog.SetupSubLoggers(Config.Logger.SubLoggers); err != nil {
		log.Fatalf("Error setting up subloggers: %s\n", err)
	}

	// Create default Version
	version := engine.CreateDefaultVersion()

	// Set build flags, unfortunately can only be of type string so must convert for IsRelease
	if BuildFlagIsRelease == "true" {
		version.IsRelease = true
		version.IsDevelopment = false
	} else {
		version.IsRelease = false
		version.IsDevelopment = true
	}
	if BuildFlagVersion != "" {
		version.Version = BuildFlagVersion
	}

	// Parse runtime flags into Version
	flag.BoolVar(&version.IsDaemon, "daemon", false, "run as a background service")
	flag.Parse()

	// Print banner and version
	qstlog.Infof(qstlog.Global, "\n"+engine.GetRandomBanner()+"\n"+version.GetVersionString(false))

	// Print logger info
	qstlog.Debugln(qstlog.Global, "Logger initialized.")
	qstlog.Debugf(qstlog.Global, "Using config dir: %s\n", Config.ConfigDir)

	// Print full path of log file name
	if strings.Contains(Config.Logger.Output, "file") {
		qstlog.Debugf(qstlog.Global, "Using log file: %s\n",
			filepath.Join(qstlog.LogPath, Config.Logger.LoggerFileConfig.FileName))
	}

	// Create the bot
	if Engine, err = engine.Create(Config, version); err != nil {
		log.Fatalf("Unable to create bot engine. Error: %s\n", err)
	}

	// Initialize the bot
	if err = Engine.Initialize(); err != nil {
		log.Fatalf("Unable to initialize bot engine. Error: %s\n", err)
	}

	// Run the bot
	if err = Engine.Run(); err != nil {
		log.Fatalf("Unable to start bot engine. Error: %s\n", err)
	}

	go func() {
		interrupt := system.WaitForInterrupt()
		s := fmt.Sprintf("Captured %v, shutdown requested.", interrupt)
		qstlog.Infoln(qstlog.Global, s)

		if !version.IsDaemon {
			systray.Quit()
		} else {
			Engine.Stop()
		}

		qstlog.Infoln(qstlog.Global, "Exiting.")
	}()

	if !version.IsDaemon {
		systray.Run(onReady, onExit)
	}

}

func onExit() {
	qstlog.Infoln(qstlog.Global, "Shutdown requested. Stopping all subsystems.")
	Engine.Stop()
	qstlog.Infoln(qstlog.Global, "Exiting.")
}

func onReady() {

	systray.SetIcon(images.Icon)
	systray.SetTitle("QuantstopTerminal")
	systray.SetTooltip("QuantstopTerminal")

	companyUrl := systray.AddMenuItem("Quantstop.com", "Quantstop.com")
	localUrl := systray.AddMenuItem("QuantstopTerminal", "Local Web App")
	dataDirUrl := systray.AddMenuItem("App Directory", "Local App Data Directory")
	quitBtn := systray.AddMenuItem("Quit", "Quit the app")

	go func() {

		for {
			select {
			case <-companyUrl.ClickedCh:
				_ = open.Run("https://quantstop.com/")
			case <-localUrl.ClickedCh:
				_ = open.Run("https://localhost")
			case <-dataDirUrl.ClickedCh:
				_ = open.Run(Config.ConfigDir)
			case <-quitBtn.ClickedCh:
				systray.Quit()
			}
		}

	}()

	// We can manipulate the systray in other goroutines
	/*go func() {
		//systray.SetTemplateIcon(icon.Data, icon.Data)
		//systray.SetTitle("QuantstopTerminal")
		//systray.SetTooltip("QuantstopTerminal")


		mChange := systray.AddMenuItem("Change Me", "Change Me")
		mChecked := systray.AddMenuItemCheckbox("Unchecked", "Check Me", true)
		mEnabled := systray.AddMenuItem("Enabled", "Enabled")
		// Sets the icon of a menu item. Only available on Mac.
		mEnabled.SetTemplateIcon(icon.Data, icon.Data)

		systray.AddMenuItem("Ignored", "Ignored")

		subMenuTop := systray.AddMenuItem("SubMenuTop", "SubMenu Test (top)")
		subMenuMiddle := subMenuTop.AddSubMenuItem("SubMenuMiddle", "SubMenu Test (middle)")
		subMenuBottom := subMenuMiddle.AddSubMenuItemCheckbox("SubMenuBottom - Toggle Panic!", "SubMenu Test (bottom) - Hide/Show Panic!", false)
		subMenuBottom2 := subMenuMiddle.AddSubMenuItem("SubMenuBottom - Panic!", "SubMenu Test (bottom)")

		mUrl := systray.AddMenuItem("Open UI", "Web App")
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

		// Sets the icon of a menu item. Only available on Mac.
		mQuit.SetIcon(icon.Data)

		systray.AddSeparator()
		mToggle := systray.AddMenuItem("Toggle", "Toggle the Quit button")
		shown := true
		toggle := func() {
			if shown {
				subMenuBottom.Check()
				subMenuBottom2.Hide()
				mQuitOrig.Hide()
				mEnabled.Hide()
				shown = false
			} else {
				subMenuBottom.Uncheck()
				subMenuBottom2.Show()
				mQuitOrig.Show()
				mEnabled.Show()
				shown = true
			}
		}

		for {
			select {
			case <-mChange.ClickedCh:
				mChange.SetTitle("I've Changed")
			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					mChecked.Uncheck()
					mChecked.SetTitle("Unchecked")
				} else {
					mChecked.Check()
					mChecked.SetTitle("Checked")
				}
			case <-mEnabled.ClickedCh:
				mEnabled.SetTitle("Disabled")
				mEnabled.Disable()
			case <-mUrl.ClickedCh:
				open.Run("https://quantstop.com/")
			case <-subMenuBottom2.ClickedCh:
				panic("panic button pressed")
			case <-subMenuBottom.ClickedCh:
				toggle()
			case <-mToggle.ClickedCh:
				toggle()
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()*/
}
