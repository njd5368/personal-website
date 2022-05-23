#!/bin/bash
yarn tailwindcss build -i site/css/tailwind-input.css -o site/css/main.css
go run main.go
