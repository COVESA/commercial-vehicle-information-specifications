#
# Makefile to generate specifications
#

.PHONY: clean all mandatory_targets json franca yaml csv ddsidl tests binary protobuf graphql ocf c overlays id jsonschema

all: clean mandatory_targets optional_targets

# All mandatory targets that shall be built and pass on each pull request for
# vehicle-signal-specification or vss-tools
mandatory_targets: clean json json-noexpand franca yaml binary csv graphql ddsidl id jsonschema apigear samm overlays

# Additional targets that shall be built by travis, but where it is not mandatory
# that the builds shall pass.
# This is typically intended for less maintainted tools that are allowed to break
# from time to time
# Can be run from e.g. travis with "make -k optional_targets || true" to continue
# even if errors occur and not do not halt travis build if errors occur
optional_targets: clean protobuf

TOOLSDIR?=./vss-tools
VSS_VERSION ?= 0.0
#COMMON_ARGS=-u ./spec/units.yaml -q ./spec/quantities.yaml  --strict
COMMON_ARGS=-u ./spec/units.yaml -q ./spec/quantities.yaml
# Default vspec root file  and validate extension if not overridden by command line input.
VSPECROOT=
VALIDATE=

json:
	vspec export json ${COMMON_ARGS} -s ${VSPECROOT} $(VALIDATE) -o cvis.json

json-noexpand:
	vspec export json ${COMMON_ARGS} --no-expand -s ${VSPECROOT} $(VALIDATE) -o cvis_noexpand.json

jsonschema:
	vspec export jsonschema ${COMMON_ARGS} -s ${VSPECROOT} $(VALIDATE) -o cvis.jsonschema

franca:
	vspec export franca --franca-vss-version $(VSS_VERSION) ${COMMON_ARGS} -s ${VSPECROOT} $(VALIDATE) -o cvis.fidl

yaml:
	vspec export yaml ${COMMON_ARGS} -s ${VSPECROOT} $(VALIDATE) -o cvis.yaml

csv:
	vspec export csv ${COMMON_ARGS} -s ${VSPECROOT} $(VALIDATE) -o cvis.csv

ddsidl:
	vspec export ddsidl ${COMMON_ARGS} -s ${VSPECROOT} $(VALIDATE) -o cvis.idl

# Verifies that supported overlay combinations are syntactically correct. The overlay files not available on CVIS.
overlays:
#	vspec export json ${COMMON_ARGS} -l overlays/profiles/motorbike.vspec -s ${VSPECROOT} $(VALIDATE) -o vss_motorbike.json
#	vspec export json ${COMMON_ARGS} -l overlays/extensions/dual_wiper_systems.vspec -s ${VSPECROOT} $(VALIDATE) -o vss_dualwiper.json
#	vspec export json ${COMMON_ARGS} -l overlays/extensions/OBD.vspec -s ${VSPECROOT} $(VALIDATE) -o vss_obd.json

tests:
	PYTHONPATH=${TOOLSDIR} pytest

binary:
	vspec export binary ${COMMON_ARGS} -s ${VSPECROOT} $(VALIDATE) -o cvis.binary

protobuf:
	vspec export protobuf ${COMMON_ARGS} -s ${VSPECROOT} $(VALIDATE) -o cvis.proto

graphql:
	vspec export graphql ${COMMON_ARGS} -s ${VSPECROOT} $(VALIDATE) -o cvis.graphql.ts


apigear:
	vspec export apigear ${COMMON_ARGS} -s ${VSPECROOT} $(VALIDATE) --output-dir apigear
	cd apigear && tar -czvf ../cvis_apigear.tar.gz * && cd ..

samm:
	vspec export samm ${COMMON_ARGS} -s ${VSPECROOT} $(VALIDATE) --target-folder samm
	cd samm && tar -czvf ../cvis_samm.tar.gz * && cd ..

id:
	vspec export id ${COMMON_ARGS} -s ${VSPECROOT} $(VALIDATE) -o cvis.vspec

clean:
	rm -f cvis.*
	rm -rf apigear
	rm -rf samm
