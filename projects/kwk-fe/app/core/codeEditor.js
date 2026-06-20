import CodeMirror, {Pos} from 'codemirror';
import {} from 'codemirror/mode/markdown/markdown.js'
import {fileName, humanDate, info, queryString, uri} from "../site-wide/util";
import str from '../stores/snpStore';

// CodeMirror.defineMode("inline-code", function() {
//     return {
//         startState: function() {return {inString: false, inTick: false};},
//         token: function(stream, state) {
//             if (stream.peek() === '`') {
//                 state.inTick = true;
//             }
//             if (state.inTick) {
//                 state.inString = !state.inString;
//                 state.inTick = false;
//                 stream.next();
//                 return "hide";
//             }
//
//             if (state.inString) {
//                 if (stream.skipTo('`')) { // Quote found on this line
//                     state.inTick = true;
//                 } else {
//                     stream.skipToEnd();    // Rest of line is string
//                 }
//                 return "code";          // Token style
//             } else {
//                 stream.skipTo('`') || stream.skipToEnd();
//                 return null;              // Unstyled token
//             }
//         }
//     };
// });

export default class CodeEditor extends Component {

    constructor(props) {
        super(props);
        info(props);
        let snp = JSON.parse(JSON.stringify(str.getSnp(uri(props.al))));
        this.state = {
            snp : snp,
        };
    }

    init() {
        if(this.cm) {
            this.dispose();
        }
        let {al} = this.props;
        let snp = str.getSnp(uri(al));
        if (snp === null) {
            return
        }
        str.currSnp = snp;
        str.similar();
        let te = document.getElementById('editor-' + uri(al));
        if (!te) return;
        this.cm = CodeMirror.fromTextArea(te, {
            mode: 'markdown',
            theme: 'dracula',
            lineNumbers: true,
            readOnly: false,
            lineWrapping: true,
            viewportMargin: 500
        });
        let that = this;
        this.cm.on('change', function (cd, ev) {
            cmChange(that, cd, ev);
        });
        this.cm.on('focus', function (cd) {
            cmFocus(that, cd);
        });
        this.cm.on('cursorActivity', function (cd, ev) {
            cursMov(that, cd, ev);
        });
        this.cm.setValue(snp.cnt || '');
        let query = queryString();
        info(query.t);
        if(query.t) {
            //this.cm.setSelection({line: 0, ch: 10}, {line:0, ch: 15});
            //this.cm.focus();
            //this.cm.setCursor({line: 1, ch: 5});
            //this.cm.addOverlay({ value : query.t })
            //this.cm.addOverlay("strings")
        }
    }

    componentDidMount() {
        info();
        this.init();
    }

    componentDidUpdate() {
        info(this.props.al);
        this.init();
    }

    componentWillUnmount() {
        info('dispose cm ' + uri(this.props.al));
        this.dispose();
    }

    dispose() {
        this.cm.toTextArea();
        this.cm = null;
    }

    render() {
        let {al} = this.props;
        if (this.state.snp == null) {
            return <div>Not found</div>
        }
        return (
            <div className="struct">
                <figure className="editor struct" style="position: relative;">
                    {/*id must have al from path as snp may not be loaded yet*/}
                    <textarea id={'editor-' + uri(al)} autoComplete="off" autoFocus="true"/>

                    <div className="editor-footer level">
                        <div id={'editor-line-' + uri(al)} className="subdued level-left"/>
                        <div className="subdued level-right">Saved {humanDate(this.state.snp.crtd)}</div>
                    </div>
                </figure>
            </div>)

    }
}

function cursMov(com, cm, ev) {
    let al = com.props.al;
    let el = document.getElementById('editor-line-' + uri(al));
    let c = cm.getCursor();
    el.innerText = c.line + ':' + c.ch;
}

function cmChange(com, cdmr, ev) {
    let cnt = cdmr.getValue();
    if(cnt === com.state.snp.cnt) {
        return;
    }
    str.currSnp.cnt = cnt;
}

function cmFocus(com, cdmr) {
    //cdmr.extendSelection(Pos(cdmr.lastLine()));
}