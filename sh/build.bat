@echo off
git describe --all --long | pipe run "in=text.cut(-8, -1)=$buildId&in.file('version.tmp')=text.replace('{buildId}', $buildId)=out.file('internal/version.go')"
go install