package openers

const JavaMain = `
import java.util.*;
import java.lang.*;
import java.io.*;
{{imports}}

class {{key}}
{
	public static void main (String[] args) throws java.lang.Exception
	{
		{{uri}}
	}
}`

const JavaImport = `import {{import}};`

const BashMain = `{{uri}}`

const NodeJsMain = `
{{imports}}

{{uri}}`

const NodeJsImport = `require('{{import}}')`

const PhpMain = `
{{imports}}

{{uri}}`

const PhpImport = `include '{{import}}'`

const PythonMain = `
{{imports}}

{{uri}}`

const PythonImport = `import {{import}}`

const GoMain = `
package main

import (
	{{imports}}
)

func main() {
 {{uri}}
}
`
const GoImport = `"{{import}}"`

// We'll not allow the 'use' statement so implementers will have to specify the fill path to modules in the main
const RustMain = `
fn main() {
    {{uri}}
}
`

const RustImport = `extern crate {{import}};`

const CSharpMain = `
{{imports}}

public class {{key}}
{
   public static void Main()
   {
      {{uri}}
   }
}`

const CSharpImport = `using {{import}};`

const ScalaMain = `
object {{key}} extends App {
   {{uri}}
 }
`
const ScalaImport = `import {{import}};`
