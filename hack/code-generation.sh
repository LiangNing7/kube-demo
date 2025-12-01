#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE[0]})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; \
    ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null \
    || echo ../code-generator)}

source "${CODEGEN_PKG}/kube_codegen.sh"

kube::codegen::gen_helpers \
    "${SCRIPT_ROOT}/pkg/apis" \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt"

kube::codegen::gen_client \
    "${SCRIPT_ROOT}/pkg/apis" \
    --with-watch \
    --output-dir "${SCRIPT_ROOT}/pkg/generated" \
    --output-pkg github.com/LiangNing7/kube-demo/pkg/generated \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt"

report_filename="${SCRIPT_ROOT}/pkg/api/api-rules/demo_violation_exceptions.list"
if [[ "${UPDATE_API_KNOWN_VIOLATIONS:-}" == "true" ]]; then
    update_report="--update-report"
fi

kube::codegen::gen_openapi \
    "${SCRIPT_ROOT}/pkg/apis" \
    --output-dir "${SCRIPT_ROOT}/pkg/generated/openapi" \
    --output-pkg github.com/LiangNing7/kube-demo/pkg/generated/openapi \
    --report-filename "${report_filename:-"/dev/null"}" \
    ${update_report:+"${update_report}"} \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt"
