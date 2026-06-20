import str from '../stores/snpStore'
import Editor from './editor';
import {info, pchPath, uri} from '../site-wide/util';


export default class EditLs extends Component {
    constructor(props) {
        super(props);
        info(props);
        let t = this;
        str.subscribe(['snp-changed'],null, (chan, val) => {
            if(str.getSelected().length>0) {
                t.context.router.push(uri(val.al));
            }
            t.setState({});
        });
    }

render() {
        info();
        return(<div className="struct editor-chrome">{
                str.getSelected().map(al => {
                    return <Editor al={al}/>;
                })
            }</div>
        )
    }
}