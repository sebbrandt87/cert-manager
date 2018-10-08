/*
Copyright 2018 The Jetstack cert-manager contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	"github.com/jetstack/cert-manager/pkg/controller"
)

func ValidateCertificateForIssuer(crt *v1alpha1.Certificate, issuerObj v1alpha1.GenericIssuer) field.ErrorList {
	el := field.ErrorList{}

	path := field.NewPath("spec")

	issuerType, err := controller.NameForIssuer(issuerObj)
	if err != nil {
		el = append(el, field.Invalid(path, err.Error(), err.Error()))
		return el
	}

	switch issuerType {
	case controller.IssuerACME:
		el = append(el, ValidateCertificateForACMEIssuer(&crt.Spec, issuerObj.GetSpec(), path)...)
	case controller.IssuerCA:
		el = append(el, ValidateCertificateForCAIssuer(&crt.Spec, issuerObj.GetSpec(), path)...)
	case controller.IssuerVault:
		el = append(el, ValidateCertificateForVaultIssuer(&crt.Spec, issuerObj.GetSpec(), path)...)
	case controller.IssuerSelfSigned:
		el = append(el, ValidateCertificateForSelfSignedIssuer(&crt.Spec, issuerObj.GetSpec(), path)...)
	}

	return el
}

func ValidateCertificateForACMEIssuer(crt *v1alpha1.CertificateSpec, issuer *v1alpha1.IssuerSpec, specPath *field.Path) field.ErrorList {
	el := field.ErrorList{}

	if crt.IsCA {
		el = append(el, field.Invalid(specPath.Child("isCA"), crt.KeyAlgorithm, "ACME does not support CA certificates"))
	}

	if len(crt.Organization) != 0 {
		el = append(el, field.Invalid(specPath.Child("organization"), crt.Organization, "ACME does not support setting the organization name"))
	}

	return el
}

func ValidateCertificateForCAIssuer(crt *v1alpha1.CertificateSpec, issuer *v1alpha1.IssuerSpec, specPath *field.Path) field.ErrorList {
	el := field.ErrorList{}

	return el
}

func ValidateCertificateForVaultIssuer(crt *v1alpha1.CertificateSpec, issuer *v1alpha1.IssuerSpec, specPath *field.Path) field.ErrorList {
	el := field.ErrorList{}

	if crt.IsCA {
		el = append(el, field.Invalid(specPath.Child("isCA"), crt.KeyAlgorithm, "Vault issuer does not currently support CA certificates"))
	}

	if len(crt.Organization) != 0 {
		el = append(el, field.Invalid(specPath.Child("organization"), crt.Organization, "Vault issuer does not currently support setting the organization name"))
	}

	return el
}

func ValidateCertificateForSelfSignedIssuer(crt *v1alpha1.CertificateSpec, issuer *v1alpha1.IssuerSpec, specPath *field.Path) field.ErrorList {
	el := field.ErrorList{}

	return el
}
