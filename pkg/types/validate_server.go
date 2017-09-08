package types

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func validateServer(s *WeblogicServer) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateServerSpec(s.Spec, field.NewPath("spec"))...)
	allErrs = append(allErrs, validateServerStatus(s.Status, field.NewPath("status"))...)
	return allErrs
}

func validateServerSpec(s WeblogicServerSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// Temporary limitation for first release.
	if s.Replicas != 1 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("replicas"), s.Replicas, "replicas can currently only be set to 1"))
	}

	allErrs = append(allErrs, validateVersion(s.Version, fldPath.Child("version"))...)

	return allErrs
}

func validateServerStatus(s WeblogicServerStatus, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validatePhase(s.Phase, fldPath.Child("phase"))...)
	return allErrs
}

func validateVersion(version string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for _, validVersion := range validVersions {
		if version == validVersion {
			return allErrs
		}
	}
	return append(allErrs, field.Invalid(fldPath, version, "invalid version specified"))
}

func validatePhase(phase WeblogicServerPhase, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for _, validPhase := range WeblogicServerValidPhases {
		if phase == validPhase {
			return allErrs
		}
	}
	return append(allErrs, field.Invalid(fldPath, phase, "invalid phase specified"))
}
