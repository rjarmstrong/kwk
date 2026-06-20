import { linkEvent } from 'inferno';
import { httpGet } from '../api/client';

export default class Snippet  extends Component {
	    constructor(props) {
		super(props);
		this.state = {
		    items : [],
		};
	    }
	
	    componentWillMount() {
		httpGet(this.props.params.u + '/' + this.props.params.p + '/' + this.props.params.n ).then(res => {
			this.setState(res || {});
		});
	    }
	
		render() {
			let s = { content : 'That snippet doesn\'t exist' };
			if (this.state.items.length>0) {
				s = this.state.items[0] || {}
			}
		    return (
		    <div>
				Transfer: {this.state.transfer}
<figure className="highlight" style="position: relative">
	<pre style="max-height: none;">
		<code className="language-html" data-lang="html" onMouseUp={linkEvent(this, snipAction)}>
		{s.content}
	</code></pre>
	<button className="button is-small bd-copy">Copy</button>
	{/* <button className="button is-small bd-expand">Expand</button> */}
	</figure>
		    </div>)
	    }
	}

	function snipAction(instance, event){
		let t = getSelectedText();
		console.log(t);
		if(!t) return;
		instance.setState((prev, props) => {
			prev.transfer = 'true';
			return prev;
		})
	}

	function getSelectedText() {
		let text = "";
		if (typeof window.getSelection !== "undefined") {
			text = window.getSelection().toString();
		} else if (typeof document.selection !== "undefined" && document.selection.type !== "Text") {
			text = document.selection.createRange().text;
		}
		return text;
	}