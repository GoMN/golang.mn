/**
 * This file/module contains all configuration for the build process.
 * src: https://gist.github.com/OverZealous/8551946
 */

/**
 * Load requires and directory resources
 */
var join = require('path').join,
  bowerrc = JSON.parse(require('fs').readFileSync('./.bowerrc', {encoding: 'utf8'})),
  bowerJSON = bowerrc.json.replace(/^\.?\/?/, './'),
  bower = require(bowerJSON),
  pkg = require('./package.json'),

  /**
   * The `buildDir` folder is where our projects are compiled during
   * development and the `compileDir` folder is where our app resides once it's
   * completely built.
   */
  buildDir = 'build',
  srcDir = 'src',
  appDir = srcDir + '/static/scripts',
  compileDir = 'dist',
  vendorDir = bowerrc.directory,
  serverTemplatesDir = 'src/templates',
  clientTemplatesDir = 'src/static/scripts/',
  indexFile = '/templates/layout.html',
  jsDir = 'src/static/scripts',
  cssDir = 'src/static/styles',
  assetsDir = 'src/static/images';

module.exports = {
    buildDir: buildDir,
    compileDir: compileDir,
    srcDir: srcDir,
    appDir: appDir,
    // Relative paths to core files and folders for input and output
    indexFile: indexFile,
    jsDir: jsDir,
    cssDir: cssDir,
    assetsDir: assetsDir,
    vendorDir: vendorDir,
    serverTemplatesDir: serverTemplatesDir,
    clientTemplatesDir: clientTemplatesDir,

    // allows settings reuse from package.json and bower.json
    bowerJSON: bowerJSON,
    bower: bower,
    pkg: pkg,

    /**
     * This code is wrapped around the application code.  It is used to "protect"
     * the application code from the global scope.
     */
    moduleWrapper: {
        header: '\n(function ( window, angular, google, $app ) {\n\n',
        footer: '\n})( window, window.angular, window.google, window.$app );'
    },

    /**
     * Settings for the server task
     * When run, this task will start a connect server on
     * your build directory, great for livereload
     */
    server: {
        port: 8080, // 0 = random port
        host: null, // null/falsy means listen to all, but will auto open localhost

        // Enable disable default auto open
        // false: run with --open to open
        // true: run with --no-open to not open, recommended if port is 0
        openByDefault: false,

        // set to false to prevent request logging
        // set to any non-`true` value to configure the logger
        log: false,

        // Live Reload server port
        lrPort: 35729
    },

    /**
     * Options passed into the various tasks.
     * These are usually passed directly into the individual gulp plugins
     */
    taskOptions: {
        csso: false, // set to true to prevent structural modifications
        jshint: {
            eqeqeq: true,
            camelcase: true,
            freeze: true,
            immed: true,
            latedef: true,
            newcap: true,
            undef: true,
            unused: true,
            browser: true,
            globals: {
                angular: false,
                console: false
            }
        },
        less: {},
        recess: {
            strictPropertyOrder: false,
            noOverqualifying: false,
            noUniversalSelectors: false
        },
        uglify: {}
    },

    /**
     * This is a collection of file patterns that refer to our app code (the
     * stuff in `src/`). These file paths are used in the configuration of
     * build tasks.
     *
     * js - All project javascript, less tests
     * jsunit - All the JS needed to run tests (in this setup, it uses the build results)
     * tpl - contains our various templates
     * html - just our main HTML file
     * less - our main stylesheet
     * assets - the rest of the files that are copied directly to the build
     */
    appFiles: {
        js: [ join(appDir, '/**/!(app)*.js'), join(appDir, '/**/*.js'), join('!', srcDir, '/**/*.spec.js'), join('!', srcDir, 'vendor/**/*.js')],
        jsunit: [ join(srcDir, '/**/*.js'), join(srcDir, '/**/*.spec.js'), '!' + join(srcDir, '/assets/**/*.js'), join('!', vendorDir, '**/*.js') ],

        tpl: [ join(appDir, '/**/*-tmpl.html') ],

        html: join(srcDir, indexFile),
        less: 'src/less/main.less',
        assets: join(srcDir, assetsDir, '**/*.*')
    },

    /**
     * Similar to above, except this is the pattern of files to watch
     * for live build and reloading.
     */
    watchFiles: {
        js: [ join(srcDir, '/script/**/*.js'), join('!', srcDir, '/**/*.spec.js'), join('!', srcDir, '/vendor/**/*.js') ],
        //jsunit: [ 'src/**/*.spec.js' ], // watch is handled by the karma plugin!

        tpl: [ join(srcDir, '/script/**/*-tmpl.html') ],

        html: [ join(buildDir, '**/*'), '!' + join(buildDir, indexFile), join(srcDir, indexFile) ],
        //less: [ 'src/**/*.less' ],
        assets: join(srcDir, assetsDir, '**/*.*')
    },

    /**
     * This is a collection of files used during testing only.
     */
    testFiles: {
        config: 'karma/karma.conf.js',
        js: [
            'vendor/angular-mocks/angular-mocks.js'
        ]
    },

    /**
     * This contains files that are provided via bower.
     * Vendor files under `js` are copied and minified, but not concatenated into the application
     * js file.
     * Vendor files under `jsConcat` are included in the application js file.
     * Vendor files under `assets` are simply copied into the assets directory.
     */
    vendorFiles: {
        js: [
            join(srcDir, '/vendor/angular/angular.js'),
            join(srcDir, '/vendor/angular-ui-router/release/angular-ui-router.js')
        ],
        jsConcat: [
            join(srcDir, '/scripts/core.js'),
            join(srcDir, 'scripts/main.js'),
            join(srcDir, 'scripts/app.js'),
            join(srcDir, 'scripts/**/*-module.js'),
            join(srcDir, 'scripts/**/*-directive.js'),
            join(srcDir, 'scripts/**/*-service.js'),
            join(srcDir, 'scripts/**/*-controller.js')
        ],
        assets: [
        ]
    },

    /**
     * This contains details about files stored on a CDN, using gulp-cdnizer.
     * file: glob or filename to match for replacement
     * package: used to look up the version info of a bower package
     * test: if provided, this will be used to fallback to the local file if the CDN fails to load
     * cdn: template for the CDN filename
     */
    cdn: [
        {
            file: 'js/vendor/angular/angular.js',
            package: 'angular',
            test: 'window.angular',
            cdn: '//ajax.googleapis.com/ajax/libs/angularjs/${ major }.${ minor }.${ patch }/angular.min.js'
        }
    ]
};