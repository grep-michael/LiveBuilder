package frontend

import (
	//appstate "LiveBuilder/AppState"
	filesystem "LiveBuilder/Filesystem"
	filelistwidgets "LiveBuilder/frontend/FileListWidgets"
	livebuildconfig "LiveBuilder/frontend/LiveBuildConfig"
	//"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/widget"
)

func buildFileSelectionView() *fyne.Container {
	package_list_widget := filelistwidgets.NewFileListContainer(filesystem.PACKAGE_DIR_ID)
	scripts_list_widget := filelistwidgets.NewFileListContainer(filesystem.CUSTOMFILES_DIR_ID)
	tabs := container.NewAppTabs(
		container.NewTabItem("Packages", package_list_widget.GetContainer()),
		container.NewTabItem("Custom Files", scripts_list_widget.GetContainer()),
	)
	return container.NewBorder(nil, nil, nil, nil, tabs)
}

func buildLBConfigView() *fyne.Container {

	return livebuildconfig.NewlbTextEdit()

}
