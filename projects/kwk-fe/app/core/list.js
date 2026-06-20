

export default class List extends Component {
    constructor(props) {
        super(props);
    }
    render({children}) {
        let {items} = this.props;
        let rows = [];
        for (let i = 0; i < items.length; i++) {
            let p = items[i];
            let nn = Inferno.cloneVNode(children, p);
            rows.push(<div>{nn}</div>)
        }
        return (<div className="scroll">{rows}</div>);
    }
}