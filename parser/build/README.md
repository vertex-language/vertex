
$ wget http://www.antlr.org/download/antlr-4.13.2-complete.jar

$ java -jar antlr-4.13.2-complete.jar -Dlanguage=Go -visitor -o ../ VertexLexer.g4 VertexParser.g4 -no-listener

$ go doc -short -all ../ > ../PACKAGE.md