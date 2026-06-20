import str from '../stores/snpStore';
import {Link} from 'inferno-router'

function getClass(score) {
    if (score >= 0.8) {
        return 'orange'
    } else if (score >= 0.7) {
        return 'yellow'
    }
}

export default class Similar extends Component{
    constructor({props}) {
        super(props);
        this.state = {
          similar : {
              cnt : [],
              time : []
          }
        };
        str.subscribe(["similar-fetched"], null, () => {
            this.setState({ similar : str.currSnp.similar })
        });
    }

    render() {
        return <div>{
            this.state.similar.cnt.map(x => {
                return <div><Link to={x.uri}>{x.uri} <span className={getClass(x.score)}>{Math.round(x.score*100)}%</span></Link> <small>merge</small></div>
            })
        }</div>
    }
}