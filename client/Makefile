compiled.js: base.js net.js client.js
	java -jar closure-compiler/compiler.jar \
		--compilation_level ADVANCED_OPTIMIZATIONS \
		--language_in ECMASCRIPT5_STRICT \
		--js base.js \
		--js net.js \
		--js client.js \
		--js_output_file compiled.js
