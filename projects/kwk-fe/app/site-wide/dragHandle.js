import {linkEvent} from 'inferno';
import {pchPath, uri, elementsFromPoint} from '../site-wide/util';
import str from "../stores/snpStore";
import {fileName} from "./util";

export default class DragHandle extends Component {

    constructor(props) {
        super(props);
        this.state = initialState();
    }

    componentDidMount() {
        window.addEventListener('mousemove', evt => {
            handleMouseMove(this, evt);
        });
        window.addEventListener('mouseup', evt => {
            handleMouseUp(this, evt);
        })
    }

    componentWillUnmount() {
        // wip:
        // window.removeEventListener('mousemove', evt => {
        //     handleMouseMove(this, evt);
        // });
        // window.removeEventListener('mouseup', evt => {
        //     handleMouseUp(this, evt);
        // })
    }

    render() {
        const {dxdy: [x, y], isPressed} = this.state;
        let scale = 1;
        let width = 30;
        if (isPressed) {
            scale = 1;
        } else {
            scale = 1;
        }
        return (
            <div style={{
                WebkitTransform: `translate3d(${x}px, ${y}px, 0) scale(${scale})`,
                transform: `translate3d(${x}px, ${y}px, 0) scale(${scale})`,
                zIndex: isPressed ? 99 : 0,
            }} className="drg"
                 onMouseDown={linkEvent(this, handleMouseDown)}
            >
                {this.props.num}
            </div>
        )
    }
}

function handleMouseDown(com, ev) {
    com.setState({
        isPressed: true,
        //dxdy: [0, 0],
        xy: [ev.clientX, ev.clientY],
    });
}

function handleMouseMove(com, ev) {
    let {xy: [x, y], hoverPouch} = com.state;
    if (ev.buttons === 1 && com.state.isPressed) {
        let dxdy = [ev.clientX - x, ev.clientY - y];

        //Pouch detection
        let el = elementsFromPoint(ev.clientX, ev.clientY) || [];
        let pchEl = el.find(x => {
            return x.className && x.className.startsWith('pouch');
        });
        // set highlight
        if (pchEl) {
            pchEl.className = 'pouch drag-over';
        }
        let nHp = pchEl ? pchEl.id : null;
        // remove highlight
        if (hoverPouch !== null && hoverPouch !== nHp) {
            document.getElementById(hoverPouch).className = 'pouch';
        }
        hoverPouch = nHp;
        com.setState({dxdy, hoverPouch});
    }
}

function initialState() {
    return {
        xy: [0, 0],
        dxdy: [0, 0],
        hoverPouch: null,
        isPressed: false
    };
}

function handleMouseUp(com, ev) {
    if(!com.state.isPressed) {
        return;
    }
    /*
        1. increment stats for target pouch.
        2. remove snips from snip list
        3. rpc to move snips
        for (let i = 0; i < this.ps.length; i++) {
        if (pchPath(this.ps[i]) === this.hoverPouch) {
            this.ps[i].stats.snips += 1;
        }
    }
     */
    if (com.state.hoverPouch !== null) {
        document.getElementById(com.state.hoverPouch).className = 'pouch';
    }
    //this.notify('drop-snippets', null, {msX: 0, msY: 0, dxdy: [0, 0]});
    com.setState(initialState());
    let sel = str.getSelected();
    sel.forEach(s => {
        let ix = str.currPch.items.findIndex(i => { return uri(i.al) === uri(s); });
        str.currPch.items.splice(ix, 1);
    });
    com.context.router.push(pchPath(str.currPch.p) + '/' + fileName(str.currPch.items[0].al));
}
