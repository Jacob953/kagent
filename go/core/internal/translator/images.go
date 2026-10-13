package translator

import (
	"regexp"

	apiv1 "github.com/kagent-dev/kagent/go/api/gen/kagent/api/v1alpha1"
)

var pinnedHarnessImage = regexp.MustCompile(`^[^[:space:]@]+@sha256:[a-f0-9]{64}$`)

func harnessImage(harness *HarnessConfiguration, images *apiv1.HarnessImagesConfig) (string, error) {
	var image string
	if override := harness.Spec.Workload.Image; override != nil {
		image = *override
	} else {
		switch harnessType(harness) {
		case HarnessTypeKagent:
			image = images.GetKagent()
		case HarnessTypeClaude:
			image = images.GetClaude()
		case HarnessTypeCodex:
			image = images.GetCodex()
		case HarnessTypeBYO:
			return "", NewValidationError("BYO harnesses must specify workload.image")
		}
		if image == "" {
			return "", NewValidationError("no default image configured for %s Harness; configure the release's built-in images or set workload.image", harnessType(harness))
		}
	}
	if !pinnedHarnessImage.MatchString(image) {
		return "", NewValidationError("Harness image must be pinned by sha256 digest")
	}
	return image, nil
}
