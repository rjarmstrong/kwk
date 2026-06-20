import str from '../stores/snpStore'
import PchItem from './pchItem'
import SnpItem from "./snpItem";

export default class Settings extends Component {
    constructor(props) {
        super(props);
        str.ps = [{n: 'zing', items: []}, {n: 'ding', items: []}];
        str.currPch.items = [{al:{n:'zong', x :'txt'}}, {al:{n:'zing', x:'xls'}}, {al:{u:'richard',p:'pipe',n:'pipe',x:'cs'}}];
        this.state = {};
    }

    render() {
        return (
            <div style="width:500px; height:500px; border: 1px solid white;">
                <div>{
                    str.ps.map(x => {
                        return <PchItem p={x}/>
                    })
                }</div>
                <div>
                    {
                        str.currPch.items.map(x => {
                            return <SnpItem s={x} />
                        })
                    }
                </div>
            </div>
        );
    };
}