for i in 1 2 3 4 5 6 7 8 9; do awk 'BEGIN { FS="\t"; OFS="," } {$1=$1; print}' The\ Vaccine\ Adverse\ Event\ Reporting\ System\ \(VAERS\)\($i\).txt |csv2json |jq . >> vaers.json; done
