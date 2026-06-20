import { httpGet } from '../api/client';

export default class FindResults  extends Component {

    constructor(props) {
        super(props);
        this.state = {
            alpha : {
                results : []
            }
        };
    }

    componentWillMount() {
        httpGet(this.props.params.u).then(res => {
            if (res.results) {
                this.setState(res);
            }
        });
    }

	render() {
	    return (
            <div>
                {JSON.stringify(this.props)}
	        <ul>{
                    this.state.alpha.results.map(x =>
                    <li>{x.highlights['main']}</li>
                    )
                }
            </ul>
            </div>
        )
    }
}