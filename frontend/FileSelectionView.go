package frontend

import (
	filesystem "LiveBuilder/Filesystem"
	filelistwidgets "LiveBuilder/frontend/FileListWidgets"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func buildFileSelectionView() *fyne.Container {
	package_list_widget := filelistwidgets.NewFileListContainer(filesystem.PACKAGE_DIR_ID)
	scripts_list_widget := filelistwidgets.NewFileListContainer(filesystem.SCRIPTS_DIR_ID)

	tabs := container.NewAppTabs(
		container.NewTabItem("Packages", package_list_widget.GetContainer()),
		container.NewTabItem("Scripts", scripts_list_widget.GetContainer()),
	)

	return container.NewBorder(nil, nil, nil, nil, tabs)

}
