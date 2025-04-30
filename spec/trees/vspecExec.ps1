# PowerShell Script

$COMMON_ARGS = "-u ../units.yaml -q ../quantities.yaml"
$VALIDATE = "-e validate"
$OVERLAY = $args[2]  # Equivalent to $3 in Bash

function Export-Yaml {
    $command = "vspec export $($args[0]) $COMMON_ARGS $OVERLAY -s $($args[1]) $VALIDATE -o exporterData/cvis.yaml"
    "vspec export $command"| Out-File -FilePath exporterData/exporterLog.txt
    Invoke-Expression $command >> "exporterData/exporterLog.txt"
}

function Export-Json {
    $command = "vspec export $($args[0]) $COMMON_ARGS $OVERLAY -s $($args[1]) $VALIDATE -o exporterData/cvis.json"
   "vspec export $command"| Out-File -FilePath exporterData/exporterLog.txt
    Invoke-Expression $command >> "exporterData/exporterLog.txt"
}

function Export-Binary {
    $command = "vspec export $($args[0]) $COMMON_ARGS $OVERLAY -s $($args[1]) $VALIDATE -o exporterData/cvis.binary"
    "vspec export $command"| Out-File -FilePath exporterData/exporterLog.txt
    Invoke-Expression $command >> "exporterData/exporterLog.txt"
}

function Export-Csv {
    $command = "vspec export $($args[0]) $COMMON_ARGS $OVERLAY -s $($args[1]) $VALIDATE -o exporterData/cvis.csv"
   "vspec export $command"| Out-File -FilePath exporterData/exporterLog.txt
    Invoke-Expression $command >> "exporterData/exporterLog.txt"
}

switch ($args[0]) {
    "yaml"   { Export-Yaml $args[0] $args[1] }
    "json"   { Export-Json $args[0] $args[1] }
    "binary" { Export-Binary $args[0] $args[1] }
    "csv"    { Export-Csv $args[0] $args[1] }
    "all" {
        Export-Yaml "yaml" $args[1]
        Export-Json "json" $args[1]
        Export-Binary "binary" $args[1]
        Export-Csv "csv" $args[1]
    }
    default {
        Write-Host "Invalid format specified. Use yaml, json, binary, csv, or all."
    }
}