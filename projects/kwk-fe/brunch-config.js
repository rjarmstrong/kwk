exports.npm = {
    globals: {
        Inferno: 'inferno',
        Component: 'inferno-component',
        CodeMirror: 'codemirror'
    },
    styles: {
        codemirror: ['lib/codemirror.css']
    }
};

exports.files = {
    javascripts: {
        entryPoints: {
            'app/app.js': 'app.js',
        }
    },
    stylesheets: {
        joinTo: 'app.css'
    }
};

exports.modules = {
    autoRequire: {
        'app.js': ['app']
    }
};

exports.plugins = {
    sass: {
        options: {
            includePaths: ["node_modules/bulma"], // tell sass-brunch where to look for files to @import
            precision: 8
        }
        // more config parameters: https://github.com/brunch/sass-brunch
    }
}