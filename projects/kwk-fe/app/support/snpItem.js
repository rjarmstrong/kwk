import {linkEvent} from 'inferno';
import str from '../stores/snpStore';
import {fileName, humanDate, info, snpNav, uri} from '../site-wide/util';


export default class SnpItem extends Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        let {al, gridView} = this.props;
        if (gridView) return;
        let el = document.getElementById(uri(al));
        el.addEventListener('contextmenu', function (e) {
            str.notify('invoke-snp-context', 'snpContext', { visible : true, al : al,  x : e.clientX, y : e.clientY });
            e.preventDefault();
        }, false);
    }

    render() {
        let s = this.props;
        if (!s) return;
        let cn = "snippet-item no-select";
        if (str.isSelected(this.props.al)) {
            cn += ' selected';
        }
        return (
            <div id={uri(s.al)} className={cn} onMouseDown={linkEvent(this, handleMouseDown)}>
                <div className="snp-wrpr">
                    <div><div className="snp-icon">✦</div>{fileName(s.al)} {s.al.draft ? <span className="draft">unsaved</span> : ''}</div>
                    { s.showPrvw === "true" ? <pre className="prvw">{s.cnt}</pre> : '' }
                    { s.showDate === "true" ? <span>{humanDate(s.crtd)}</span> : '' }
                </div>
            </div>
        );
    };
}

function handleMouseDown(com, ev) {
    let {al, gridView} = com.props;
    if (ev.buttons !== 1 && ev.buttons !== 3) {
        return;
    }
    // temporary while we think about multi-select
    // in gridView
    if (ev.buttons !== 1 && gridView) {
        return;
    }
    let sel = str.getSelected();
    // if right click, multiple items selected and
    // one is the click target then exit func.
    if (ev.buttons === 3 && sel.length > 1) {
        let found = sel.find(x => { return uri(al) === uri(x); });
        if (found !== null) return;

    }
    let dest = [al];
    if (ev.metaKey) {
        dest = str.recreateMal(al, sel);
    }
    snpNav(com.context.router, dest);
}