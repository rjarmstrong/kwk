import {info} from '../site-wide/util';
import SignUp from '../site-wide/signUp'

export default class PubRoot extends Component {

    constructor(props) {
        info();
        super(props);
        this.state = {}
    }

    render() {
        return (<div>
                <div className="columns">
                    <div className="column is-6">
                        <section className="hero">
                            <div className="hero-body">
                                <div className="container">
                                    <h1 className="title">
                                        Little notes made big.
                                    </h1>
                                    <h2 className="subtitle">
                                        Stop scratching around and get things done.
                                    </h2>
                                    <p>kwk is a friendly and dedicated app to get in control of your notes<br/> so you can spend more time on your craft.<br/>
                                    </p>
                                    <br/>
                                    <h3>Key Features</h3>
                                    <ul>
                                        <li>- Lightweight and quick note taking</li>
                                        <li>- A streamlined search built specially for notes.</li>
                                        <li>- De-duplication tools to stop repeating yourself</li>
                                        <li>- Merge related concepts</li>
                                        <li>- Easy to remember links</li>
                                    </ul>
                                </div>
                            </div>
                        </section>
                    </div>
                    <div className="column is-4">
                        <h3>Sign Up</h3>
                        {/*<p>Stop scratching around and get productive.</p>*/}
                        <SignUp/>
                    </div>
                </div>
                {/*<div class="container is-widescreen">*/}
                {/*<div class="notification">*/}
                {/*This container is <strong>fluid</strong>: it will have a 24px gap on either side, on any viewport size.*/}
                {/*</div>*/}
                {/*</div>*/}

            </div>
        )
        info();
    }
}