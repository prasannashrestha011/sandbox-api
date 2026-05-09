package sandbox_image

func LoadImage(env string) string {
	switch env {

	case "python":
		return "python:3.11-slim"

	case "java":
		return "eclipse-temurin:21-jdk-jammy"

	case "javascript":
		return "node:20-alpine"

	case "c":
		return "gcc:alpine"

	case "cpp":
		return "gcc:alpine"

	default:
		return ""
	}
}
