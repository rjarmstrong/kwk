package execution

import (
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"os/exec"
	"testing"
	//"os"
	"bytes"
	//"io"
	//"fmt"
	//"os"
	//"io"
)

func Test_Execution(t *testing.T) {
	Convey("Exec safe", t, func() {

		Convey(`Should pipe execution across multiple functions`, func() {
			script := `
			 process.stdin.on('readable', () => {
			  var chunk = process.stdin.read();
			  if (chunk !== null) {
			    process.stdout.write(chunk.toString());
			  }
			});

			//process.stdin.on('end', () => {
			//  process.stdout.write('end');
			//});
			`
			var b bytes.Buffer
			err := Execute(&b,
				exec.Command("node", "-e", `
					process.stdout.write('hola'+ ", ");
					process.stdout.write('zingo'+ ", ");
				`),
				exec.Command("node", "-e", script),
				exec.Command("/bin/bash", "-c", `echo "start| $(cat -)|end"`),
			)
			So(err, ShouldBeNil)
			So(b.String(), should.Equal, "start| hola, zingo, |end\n")
			//io.Copy(os.Stdout, &b)
		})

		Convey(`Should pass arguments and then pipe to multiple functions`, func() {
			script := `
				for(var i = 1; i < process.argv.length; i++) {
  					process.stdout.write(process.argv[i] + ", ");
				};
			`
			var b bytes.Buffer
			err := Execute(&b,
				exec.Command("node", "-e", script, "arg1", "arg2"),
				exec.Command("/bin/bash", "-c", `echo "start| $(cat -)|end"`),
			)
			So(err, ShouldBeNil)
			So(b.String(), should.Equal, "start| arg1, arg2, |end\n")
			//io.Copy(os.Stdout, &b)
		})

	})
}
