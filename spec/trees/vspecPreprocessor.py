import json
import os
from argparse import ArgumentParser

def generate_2dim_overlay(inst_matrix, instance_name, prefix_data, output):
    num_dim1 = len(inst_matrix)

    # Dimension 1 instance line
    output.append(prefix_data[0] + ":")
    dim2_names = []
    for dim1member in inst_matrix[0]:
        dim2_names.append(dim1member)
    output.append(f"  instances: {list(dim2_names)}")
    output.append("")

    # Dim2 for each dim1 member
    for dim1_index, dim1member_name in enumerate(inst_matrix[0]):
        dim2 = inst_matrix[1][dim1_index]
        output.append(f"{prefix_data[0]}.{dim1member_name}{prefix_data[1]}:")
        output.append(f"  instances: {dim2}")
        output.append("")

    return output

def generate_3dim_overlay(inst_matrix, prefix_data, output):
    num_dim1 = len(inst_matrix)

    # Dimension 1 instance line
    output.append(prefix_data[0] + ":")
    dim2_names = []
    for dim1member in inst_matrix:
        dim2_names.append(list(dim1member.keys())[0])
    output.append(f"  instances: {list(dim2_names)}")
    output.append("")

    # Iterate over dim1 members
    for dim1member in inst_matrix:
        dim1member_name = list(dim1member.keys())[0]
        dim1member_data = dim1member[dim1member_name][list(dim1member[dim1member_name].keys())[0]]
        dim2 = dim1member_data[0]
        dim3 = dim1member_data[1]

        # Dim2 section
        output.append(f"{prefix_data[0]}.{dim1member_name}{prefix_data[1]}:")
        output.append(f"  instances: {dim2}")
        output.append("")

        # Dim3 for each dim2 member
        for dim2_index, dim2member_name in enumerate(dim2):
            dim3member = dim3[dim2_index]
            output.append(f"{prefix_data[0]}.{dim1member_name}{prefix_data[1]}.{dim2member_name}{prefix_data[2]}:")
            output.append(f"  instances: {dim3member}")
            output.append("")

    return output

def get_instance_prefix_data(instance_name1, scope_data):
    prefix = []

    config_scope = scope_data["config-scope"]
    for cfg_scope in config_scope:
        if list(cfg_scope.keys())[0] == "instance-scope":
            prefixes = cfg_scope["instance-scope"]
            for instance in prefixes:
                instance_name2 = list(instance.keys())[0]
                if instance_name2 == instance_name1:
                    prefix = instance[instance_name2]

    return prefix

def generate_variant_overlay(variant_type_name, variant_name, scope_data, output):
    config_scope = scope_data["config-scope"]
    for cfg_scope in config_scope:
        if list(cfg_scope.keys())[0] == "variant-scope":
            variant_type_data = cfg_scope["variant-scope"]
            for variant_type in variant_type_data:
                if list(variant_type.keys())[0] == variant_type_name:
                    for variant_data in variant_type[variant_type_name]:
                        if variant_data == variant_name:
                            for overlay_data in variant_type[variant_type_name][variant_data]:
                                output.append(f'{overlay_data["Path"]}:')
                                output.append(f'{overlay_data["Directive"]}')
                                output.append("")
    return output

def get_instances(instance_name):
    return ["Trailer1", "Trailer2", "Trailer3"] #TODO

def generate_instance_variant_overlay(instance_variant_type_name, instance_variant_config_data, scope_data, output):
    config_scope = scope_data["config-scope"]
    for cfg_scope in config_scope:
        if list(cfg_scope.keys())[0] == "instance-variant-scope":
            instance_variant_type_data = cfg_scope["instance-variant-scope"]
            for instance_variant_type in instance_variant_type_data:
                if len(instance_variant_config_data) == 2:
                    instance_variant_config_name = instance_variant_config_data[0]
                    instance_variant_config_x1 = [""] # len=1
                    instance_variant_config_x2 = instance_variant_config_data[1]
                else:
                    instance_variant_config_name = instance_variant_config_data[0]
                    instance_variant_config_x1 = get_instances(instance_variant_config_data[1])
                    instance_variant_config_x2 = instance_variant_config_data[2]
                if list(instance_variant_type.keys())[0] == instance_variant_type_name:
                    for instance_variant_data in instance_variant_type[instance_variant_type_name]:
                        if instance_variant_data == instance_variant_config_name:
                            for x1 in instance_variant_config_x1:
                                for overlay_data in instance_variant_type[instance_variant_type_name][instance_variant_config_name]:
                                    output.append(f'{overlay_data["Path"].replace(".X1.", "." + x1 + ".").replace(".X2.", "." + instance_variant_config_x2 + ".")}:')
                                    output.append(f'{overlay_data["Directive"]}')
                                    output.append("")
    return output

def generate_overlay(config_data, scope_data):
    output = []

    configs = config_data["configurations"]
    for cfg in configs:
        if list(cfg.keys())[0] == "instances":
            instances = cfg["instances"]
            # Iterate over each instance configuration
            for instance in instances:
                instance_name = list(instance.keys())[0]
                instance_matrix = instance[instance_name]
                instance_prefix = get_instance_prefix_data(instance_name, scope_data)
                if len(instance_prefix) == 3:
                    output = generate_3dim_overlay(instance_matrix, instance_prefix, output)
                elif len(instance_prefix) == 2:
                    output = generate_2dim_overlay(instance_matrix, instance_name, instance_prefix, output)
                else:
                    print(f"Scope data matching {instance_name} not found\n")

        if list(cfg.keys())[0] == "variants":
            variants = cfg["variants"]
            # Iterate over each variant configuration
            for variant in variants:
                variant_type = list(variant.keys())[0]
                variant_name = variant[variant_type]
                output = generate_variant_overlay(variant_type, variant_name, scope_data, output)

        if list(cfg.keys())[0] == "instance-variants":
            instance_variants = cfg["instance-variants"]
            # Iterate over each instance-variant configuration
            for instance_variant in instance_variants:
                instance_variant_type = list(instance_variant.keys())[0]
                instance_variant_data = instance_variant[instance_variant_type]
                output = generate_instance_variant_overlay(instance_variant_type, instance_variant_data, scope_data, output)

    return "\n".join(output)

# ------------------ Load JSON config ------------------
parser = ArgumentParser(prog='python3 vspecPreprocessor.py',
                    description='The VspecPreprocessor tool configures the vspec files by creating overlays to be submitted with the vspec files to VSS-tools')
parser.add_argument("-i", "--inputfile", help="JSON file containing configuration data", required=True)
parser.add_argument("-s", "--scopefile", help="JSON file defining the scope of the different config features", default='configScope.json', required=False)
parser.add_argument("-o", "--outputfile", help="Overlay vspec file to be used with VSS-tools", default='overlay.vspec', required=False)
parser.add_argument("-v", "--vspecfile", help="Root vspec file of the tree", required=False)
parser.add_argument("-f", "--format", help="Exporter output format", choices=['yaml', 'json', 'binary'], default='yaml',required=False)

args = parser.parse_args()

input_file = args.inputfile
with open(input_file, "r") as f1:
    config_data = json.load(f1)

scope_file = args.scopefile
with open(scope_file, "r") as f2:
    scope_data = json.load(f2)

overlay_text = generate_overlay(config_data, scope_data)

#print(overlay_text)

output_file = args.outputfile
with open(output_file, "w") as f:
    f.write(overlay_text)

print(f"\nOverlay configuration saved to {output_file}")
if args.vspecfile :
    path = os.path.dirname(args.vspecfile)
    if len(path) == 0:
        print(f"\nExporter command: vspec export {args.format} -u {path}units.yaml -q {path}quantities.yaml -l {output_file} -s {args.vspecfile} -o cvis.{args.format}")
    else:
        print(f"\nExporter command: vspec export {args.format} -u {path}/units.yaml -q {path}/quantities.yaml -l {output_file} -s {args.vspecfile} -o cvis.{args.format}")
