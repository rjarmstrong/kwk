

export default class Grid extends Component {
    constructor(props) {
        super(props);
    }
    render({children}) {
        let {items} = this.props;
        let rows = [[]];
        for (let i = 0; i < items.length; i++) {
            let p = items[i];
            let row = Math.floor(i/5);
            if (rows.length-1 < row) {
                rows.push([])
            }
            let nn = Inferno.cloneVNode(children, p);
            rows[row].push(<div className="column is-2">{nn}</div>)
        }
        let grid = [];
        for (let i = 0; i < rows.length; i++) {
            grid.push(<div className="columns is-12">{rows[i]}</div>)
        }
        return (<div className="scroll">{grid}</div>);
    }
}