package media

import (
	"fmt"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

type ImageRefController struct {
	view.SetModelValueControllerBase
}

func (self ImageRefController) Supports(metaData *model.MetaData, form *view.Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*ImageRef)
	return ok
}

func (self ImageRefController) NewInput(withLabel bool, metaData *model.MetaData, form *view.Form) (input view.View, err error) {
	imageRef := metaData.Value.Addr().Interface().(*ImageRef)
	thumbnailSize := Config.ImageRefController.ThumbnailSize

	var img view.View
	if imageRef.IsEmpty() {
		img = view.IMG(Config.dummyImageURL, thumbnailSize, thumbnailSize)
	} else {
		image, err := imageRef.Get()
		if err != nil {
			return nil, err
		}
		version, err := image.Thumbnail(thumbnailSize)
		if err != nil {
			return nil, err
		}
		img = version.ViewImage("")
	}

	hiddenInput := &view.HiddenInput{Name: metaData.Selector(), Value: imageRef.String()}

	thumbnailFrame := &view.Div{
		Class:   Config.ImageRefController.ThumbnailFrameClass,
		Style:   fmt.Sprintf("width:%dpx;height:%dpx;", thumbnailSize, thumbnailSize),
		Content: img,
	}

	removeButton := &view.Button{
		Content: view.HTML("Remove"),
		OnClick: fmt.Sprintf(
			`jQuery("#%s").attr("value", "");
			jQuery("#%s").empty().append("<img src='%s' width='%d' height='%d'/>");`,
			hiddenInput.ID(),
			thumbnailFrame.ID(),
			Config.dummyImageURL,
			thumbnailSize,
			thumbnailSize,
		),
	}

	chooseDialogThumbnails := view.DIV("")
	chooseDialogThumbnailsID := chooseDialogThumbnails.ID()
	chooseDialog := &view.ModalDialog{
		Style: "width:600px;height:400px;",
		Content: view.Views{
			view.H3("Choose Image:"),
			chooseDialogThumbnails,
			view.ModalDialogCloseButton("Close"),
		},
	}

	editor := view.DIV(Config.ImageRefController.Class,
		hiddenInput,
		thumbnailFrame,
		view.DIV(Config.ImageRefController.ActionsClass,
			// view.HTML("&larr; drag &amp; drop files here"),
			chooseDialog,
			view.DynamicView(
				func(ctx *view.Context) (view.View, error) {
					ctx.Response.RequireScriptURL("/media/media.js", 0)
					return &view.Button{
						Content: view.HTML("Choose"),
						OnClick: fmt.Sprintf(
							`gostart_media.fillChooser('#%s', '%s', function(value){
								jQuery('#%s').attr('value', value.id);
								var img = '<img src=\"'+value.url+'\" alt=\"'+value.title+'\"/>';
								jQuery('#%s').empty().append(img);
								%s
							});
							%s;`,
							chooseDialogThumbnailsID,
							AllThumbnailsAPI.URL(ctx.ForURLArgsConvert(Config.ImageRefController.ThumbnailSize)),
							hiddenInput.ID(),
							thumbnailFrame.ID(),
							view.ModalDialogCloseScript,
							chooseDialog.OpenScript(),
						),
					}, nil
				},
			),
			removeButton,
			UploadImageButton(
				thumbnailSize,
				fmt.Sprintf(
					`function(id, fileName, responseJSON) {
						var img = "<img src='" + responseJSON.thumbnailURL + "' width='%d' height='%d'/>";
						jQuery("#%s").empty().html(img);
						jQuery("#%s").attr("value", responseJSON.imageID);
					}`,
					thumbnailSize,
					thumbnailSize,
					thumbnailFrame.ID(),
					hiddenInput.ID(),
				),
			),
		),
		view.DivClearBoth(),
	)

	if withLabel {
		return view.AddStandardLabel(form, editor, metaData), nil
	}
	return editor, nil
}
