package routes

import (
	"github.com/go-chi/chi/v5"

	controllers "main/internal/controllers/lab"
	"main/internal/proxy"
	"main/internal/response"
)

// RegisterLabRoutes wires lab endpoints into the router.
func RegisterLabRoutes(r chi.Router, labController *controllers.LabController, chapterController *controllers.ChapterController,
	exerciseController *controllers.ExerciseHandler, enrollmentController *controllers.EnrollmentController) {
	r.Route("/labs", func(sr chi.Router) {
		sr.Use(proxy.AuthMiddleware)
		//instructors routes
		sr.Route("/", func(lr chi.Router) {
			lr.Use(proxy.UserTypeMiddlware)
			response.WrapPost(lr, "/", labController.CreateLab)
			response.WrapGet(lr, "/{labId}", labController.GetLabByID)
			response.WrapDelete(lr, "/{labId}", labController.DeleteLab)

			lr.Route("/{labId}/chapters", func(cr chi.Router) {
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
		//students routes
		sr.Route("/enrollments", func(er chi.Router) {
			response.WrapPost(er, "/", enrollmentController.EnrollUserToLab)
			response.WrapGet(er, "/{labId}", enrollmentController.GetEnrollment)
			response.WrapGet(er, "/user/{userId}", enrollmentController.GetUserEnrollments)
			response.WrapDelete(er, "/{labId}", enrollmentController.DeleteEnrollment)
		})

	})
}
