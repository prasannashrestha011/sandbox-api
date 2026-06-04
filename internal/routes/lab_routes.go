package routes

import (
	"github.com/go-chi/chi/v5"

	controllers "main/internal/controllers/lab"
	"main/internal/proxy"
	"main/internal/response"
)

// RegisterLabRoutes wires lab endpoints into the router.
func RegisterLabRoutes(r chi.Router, labController *controllers.LabController, chapterController *controllers.ChapterController, exerciseController *controllers.ExerciseHandler) {
	r.Route("/labs", func(sr chi.Router) {
		sr.Use(proxy.AuthMiddleware)
		sr.Use(proxy.UserTypeMiddlware)
		response.WrapPost(sr, "/", labController.CreateLab)
		response.WrapGet(sr, "/{labId}", labController.GetLabByID)
		response.WrapDelete(sr, "/{labId}", labController.DeleteLab)

		sr.Route("/{labId}/chapters", func(cr chi.Router) {
			response.WrapPost(cr, "/", chapterController.CreateChapter)
			response.WrapGet(cr, "/", chapterController.GetChaptersByLabID)
			response.WrapPut(cr, "/{chapterId}", chapterController.UpdateChapter)

			cr.Route("/{chapterId}/exercises", func(er chi.Router) {
				response.WrapPost(er, "/", exerciseController.CreateExercise)
				response.WrapGet(er, "/", exerciseController.ListExercisesByChapterID)
				response.WrapGet(er, "/{exerciseId}", exerciseController.GetExerciseByID)
				response.WrapPut(er, "/{exerciseId}", exerciseController.UpdateExercise)
				response.WrapDelete(er, "/{exerciseId}", exerciseController.DeleteExercise)
			})
		})
	})
}
