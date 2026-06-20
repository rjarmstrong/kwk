#!/bin/sh

# BUILD FOR LINUX ONLY

echo "** MINIFY **"
echo "READING JS..."
DIR=$(dirname $0)
OUTDIR="/usr/share/nginx/html/"

for filename in ${OUTDIR}*.js; do
     cat "${filename}" | "${DIR}/minify" > "${filename}.tmp"
done
echo "WRITING JS..."
for filename in ${OUTDIR}*.js; do
    rm "${filename}"
done
echo "CLEANING..."
for filename in ${OUTDIR}*.js.tmp; do
     mv "${filename}" "${OUTDIR}"`basename ${filename} .tmp`
done
echo "DONE"