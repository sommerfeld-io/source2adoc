package io.sommerfeld.source2adoc

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.shell.standard.ShellMethod

@SpringBootApplication
class Source2adocApplication

fun main(args: Array<String>) {
	runApplication<Source2adocApplication>(*args)
}