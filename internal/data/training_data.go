package data

// FrameworkExamples maps framework names to code snippet examples for ML/keyword detection
var FrameworkExamples = map[string][]string{
	"React": {
		"import React from 'react'",
		"useState(",
		"useEffect(",
		"export default function",
		"<div>",
	},
	"Gin": {
		"github.com/gin-gonic/gin",
		"gin.Default()",
		"r := gin.New()",
		"r.GET(",
		"c.JSON(",
	},
	"Django": {
		"from django",
		"def get_queryset(self):",
		"class Meta:",
		"urlpatterns = [",
		"from rest_framework",
	},
	"Flask": {
		"from flask",
		"Flask(__name__)",
		"@app.route(",
		"app.run(",
	},
	"Spring": {
		"org.springframework",
		"@Controller",
		"@RequestMapping",
		"SpringApplication.run",
	},
	"Express": {
		"require('express')",
		"const app = express()",
		"app.get(",
		"app.listen(",
	},
	"Vue": {
		"import Vue from 'vue'",
		"new Vue({",
		"<template>",
	},
	"Angular": {
		"@Component({",
		"import { Component }",
		"ngOnInit()",
	},
}
