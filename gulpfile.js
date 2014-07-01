//<editor-fold desc="Node Requires, gulp, etc">
//https://gist.github.com/OverZealous/8551946
var gulp = require('gulp'),
  autoprefixer = require('gulp-autoprefixer'),
  bump = require('gulp-bump'),
  cdnizer = require('gulp-cdnizer'),
  clean = require('gulp-clean'),
  concat = require('gulp-concat'),
  csso = require('gulp-csso'),
  debug = require('gulp-debug'),
  footer = require('gulp-footer'),
  gutil = require('gulp-util'),
  gzip = require('gulp-gzip'),
  header = require('gulp-header'),
  help = require('gulp-task-listing'),
  _if = require('gulp-if'),
  inject = require('gulp-inject'),
  jshint = require('gulp-jshint'),
  karma = require('gulp-karma'),
  less = require('gulp-less'),
  livereload = require('gulp-livereload'),
  livereloadEmbed = require('gulp-embedlr'),
  minifyHtml = require('gulp-minify-html'),
  ngHtml2js = require('gulp-ng-html2js'),
  ngmin = require('gulp-ngmin'),
  plumber = require('gulp-plumber'),
  recess = require('gulp-recess'),
  rename = require('gulp-rename'),
  rev = require('gulp-rev'),
  tap = require('gulp-tap'),
  uglify = require('gulp-uglify'),
  watch = require('gulp-watch'),

  _ = require('lodash'),
  anysort = require('anysort'),
  connect = require('connect'),
  es = require('event-stream'),
  fs = require('fs'),
  http = require('http'),
  lazypipe = require('lazypipe'),
  open = require('open'),
  path = require('path'),
  runSequence = require('run-sequence'),
  server = require('tiny-lr')();
//</editor-fold>

/*
 * TODO:
 *  - Banner for compressed JS/CSS
 *  - Changelog?
 */

// Load user config and package data
var cfg = require('./build.config.js');

// common variables
var concatName = cfg.pkg.name,
  embedLR = false,
  join = path.join;



//=============================================
// MAIN TASKS
//=============================================

gulp.task('default', ['watch']);

gulp.task('server-dist', ['compile'], function(cb) {
    startServer(cfg.compileDir, cb);
});

gulp.task('server', ['watch'], function(cb) {
    startServer(cfg.buildDir, cb);
});

gulp.task('build', function(cb) {
    runSequence(['build-assets', 'build-scripts', 'build-styles'],
      ['build-html'
      //  , 'test'
    ], cb);
});
gulp.task('watch', ['lr-server', 'build', 'test-watch'], function() {
    watch({glob: cfg.watchFiles.js, emitOnGlob: false, name: 'JS'})
      .pipe(plumber())
      .pipe(jsBuildTasks())
      .pipe(livereload(server));

    watch({glob: cfg.watchFiles.tpl, emitOnGlob: false, name: 'Templates'})
      .pipe(plumber())
      .pipe(tplBuildTasks())
      .pipe(livereload(server));

    watch({glob: cfg.watchFiles.html, emitOnGlob: false, name: 'HTML'}, function() {
        return buildHTML();
    });

    watch({glob: cfg.watchFiles.less, emitOnGlob: false, name: 'Styles'}, function() {
        // run this way to ensure that a failed pipe doesn't break the watcher.
        return buildStyles();
    });

    watch({glob: cfg.watchFiles.assets, emitOnGlob: false, name: 'Assets'})
      .pipe(gulp.dest(join(cfg.buildDir, cfg.assetsDir)));
});

gulp.task('compile', function(cb) {
    runSequence(
      'compile-clean',
      ['compile-assets', 'compile-scripts', 'compile-styles'],
      'compile-html',
      cb
    );
});

gulp.task('clean', ['compile-clean', 'build-clean']);

gulp.task('help', help);




//=============================================
// UTILITIES
//=============================================

function readFile(filename) {
    return fs.existsSync(filename) ? fs.readFileSync(filename, {encoding: 'utf8'}) : '';
}

function insertRevGlob(f) {
    // allows for rev'd filenames (which end with /-[0-9a-f]{8}.ext/
    return f.replace(/(\..*)$/, '?(-????????)$1');
}

function startServer(root, cb) {
    var devApp, devServer, devAddress, devHost, url, log=gutil.log, colors=gutil.colors;

    devApp = connect();
    if(cfg.server.log) {
        devApp.use(connect.logger(cfg.server.log===true ? 'dev' : cfg.server.log));
    }
    devApp.use(connect.static(root));

    devServer = http.createServer(devApp).listen(cfg.server.port, cfg.server.host||undefined);

    devServer.on('error', function(error) {
        log(colors.underline(colors.red('ERROR'))+' Unable to start server!');
        cb(error);
    });

    devServer.on('listening', function() {
        devAddress = devServer.address();
        devHost = devAddress.address === '0.0.0.0' ? 'localhost' : devAddress.address;
        url = 'http://' + devHost + ':' + devAddress.port + join('/', cfg.indexFile);

        log('');
        log('Started dev server at '+colors.magenta(url));
        var openByDefault = cfg.server.openByDefault;
        if(gutil.env.open || (openByDefault && gutil.env.open !== false)) {
            log('Opening dev server URL in browser');
            if(openByDefault) {
                log(colors.gray('(Run with --no-open to prevent automatically opening URL)'));
            }
            // Open the URL in the browser at this point.
            open(url);
        } else if(!openByDefault) {
            log(colors.gray('(Run with --open to automatically open URL on startup)'));
        }
        log('');
        cb();
    });
}

//=============================================
// SUB TASKS
//=============================================


//---------------------------------------------
// HTML
//---------------------------------------------

var buildHTML = function() {
    var htmlFile = readFile(join(cfg.srcDir, cfg.indexFile));
    return gulp.src([join(cfg.srcDir, '/**/*.*'), '!' + join(cfg.srcDir, cfg.indexFile)], {read: false})
      .pipe(plumber())
      .pipe(inject(cfg.appFiles.html, {
          addRootSlash: false,
          sort: fileSorter, // see below
          ignorePath: join('/',cfg.buildDir,'/')
      }))
      .pipe(_if(embedLR, livereloadEmbed({port: cfg.server.lrPort})))
      .pipe(gulp.dest(cfg.buildDir))
      .pipe(tap(function(file) {
          var newHtmlFile = file.contents.toString();
          if(newHtmlFile !== htmlFile) {
              htmlFile = newHtmlFile;
              //gulp.src(file.path).pipe(livereload(server));
          }
      }));
};

gulp.task('build-html', function() {
    // NOTE: this task does NOT depend on buildScripts and buildStyles,
    // therefore, it may incorrectly create the HTML file if called
    // directly.
    return buildHTML();
});

gulp.task('compile-html', function() {
    // NOTE: this task does NOT depend on compileScripts and compileStyles,
    // therefore, it may incorrectly create the HTML file if called
    // directly.
    return gulp.src([join(cfg.compileDir, '/**/*.*'), '!' + join(cfg.compileDir, cfg.indexFile)], {read: false})
      .pipe(inject(cfg.appFiles.html, {
          addRootSlash: false,
          sort: fileSorter, // see below
          ignorePath: join('/', cfg.compileDir, '/')
      }))
      .pipe(cdnizer(cfg.cdn))
      .pipe(minifyHtml({empty:true,spare:true,quotes:true}))
      .pipe(gulp.dest(cfg.compileDir))
      .pipe(gzip())
      .pipe(gulp.dest(cfg.compileDir))
});

// used by build-html to ensure correct file order during builds
var fileSorter = (function(){
    var as = anysort(_.flatten([
        // JS files are sorted by original vendor order, common, app, then everything else
        cfg.vendorFiles.js.map(function(f){ return join(cfg.jsDir, insertRevGlob(f)); }),
        cfg.vendorFiles.jsConcat.map(function(f){ return join(cfg.jsDir, insertRevGlob(f)); }),
        join(cfg.jsDir, 'core.js'),
        join(cfg.jsDir, 'main.js'),
        join(cfg.jsDir, 'app.js'), // exclude app.js so it comes last
        join(cfg.jsDir, '**/*.js'),
        // CSS order should be maintained via Less includes
        join(cfg.cssDir, '**/*.css')
    ]));
    return function(a,b){ return as(a.filepath, b.filepath) || a.filepath < b.filepath };
})();



//---------------------------------------------
// JavaScript
//---------------------------------------------

var jsFiles = function() { return gulp.src(cfg.appFiles.js); },
  jsBaseTasks = lazypipe()
    .pipe(plumber)  // jshint won't render parse errors without plumber
    .pipe(function() {
        return jshint(_.clone(cfg.taskOptions.jshint));
    })
    .pipe(jshint.reporter, 'jshint-stylish'),
  jsBuildTasks = jsBaseTasks
    .pipe(gulp.dest, join(cfg.buildDir, cfg.jsDir)),
  tplFiles = function() { return gulp.src(cfg.appFiles.tpl); },
  tplBuildTasks = lazypipe()
    .pipe(ngHtml2js, {moduleName: 'templates'})
    .pipe(gulp.dest, join(cfg.buildDir, cfg.jsDir, cfg.clientTemplatesDir));

//noinspection FunctionWithInconsistentReturnsJS
gulp.task('build-scripts-vendor', function() {
    var vendorFiles = [].concat(cfg.vendorFiles.js||[], cfg.vendorFiles.jsConcat||[]);
    if(vendorFiles.length) {
        return gulp.src(vendorFiles, {base: cfg.vendorDir})
          .pipe(gulp.dest(join(cfg.buildDir, cfg.jsDir, cfg.vendorDir)))
    }
});
gulp.task('build-scripts-app', function() {
    return jsFiles().pipe(jsBuildTasks());
});
gulp.task('build-scripts-templates', function() {
    return tplFiles().pipe(tplBuildTasks());
});
gulp.task('build-scripts', ['build-scripts-vendor', 'build-scripts-app', 'build-scripts-templates']);


gulp.task('compile-scripts', function() {
    var appFiles, templates, files, concatFiles, vendorFiles;
    appFiles = jsFiles()
      .pipe(jsBaseTasks())
      .pipe(concat('appFiles.js')) // not used
      .pipe(ngmin())
      .pipe(header(cfg.moduleWrapper.header))
      .pipe(footer(cfg.moduleWrapper.footer));

    templates = tplFiles()
      .pipe(minifyHtml({empty: true, spare: true, quotes: true}))
      .pipe(ngHtml2js({moduleName: 'templates'}))
      .pipe(concat('templates.min.js')); // not used

    files = [appFiles, templates];
    if(cfg.vendorFiles.jsConcat.length) {
        files.unshift(gulp.src(cfg.vendorFiles.jsConcat));
    }

    concatFiles = es.concat.apply(es, files)
      .pipe(concat(concatName + '.js'))
      .pipe(uglify(cfg.taskOptions.uglify))
      .pipe(rev())
      .pipe(gulp.dest(join(cfg.compileDir, cfg.jsDir)))
      .pipe(gzip())
      .pipe(gulp.dest(join(cfg.compileDir, cfg.jsDir)));

    if(cfg.vendorFiles.js.length) {
        vendorFiles = gulp.src(cfg.vendorFiles.js, {base: cfg.vendorDir})
          .pipe(uglify(cfg.taskOptions.uglify))
          .pipe(rev())
          .pipe(gulp.dest(join(cfg.compileDir, cfg.jsDir, cfg.vendorDir)))
          .pipe(gzip())
          .pipe(gulp.dest(join(cfg.compileDir, cfg.jsDir, cfg.vendorDir)));

        return es.concat(vendorFiles, concatFiles);
    } else {
        return concatFiles;
    }
});



//---------------------------------------------
// Less / CSS Styles
//---------------------------------------------

var styleFiles = function() { return gulp.src(cfg.appFiles.less); },
  styleBaseTasks = lazypipe()
    .pipe(recess, cfg.taskOptions.recess)
    .pipe(less, cfg.taskOptions.less)
    .pipe(autoprefixer),
  buildStyles = function() {
      return styleFiles()
        .pipe(
        styleBaseTasks()
            // need to manually catch errors on recess/less :-(
          .on('error', function() {
              gutil.log(gutil.colors.red('Error')+' processing Less files.');
          })
      )
        .pipe(gulp.dest(join(cfg.buildDir, cfg.cssDir)))
        .pipe(livereload(server))
  };

gulp.task('build-styles', function() {
    return buildStyles();
});
gulp.task('compile-styles', function() {
    return styleFiles()
      .pipe(styleBaseTasks())
      .pipe(rename(concatName + '.css'))
      .pipe(csso(cfg.taskOptions.csso))
      .pipe(rev())
      .pipe(gulp.dest(join(cfg.compileDir, cfg.cssDir)))
      .pipe(gzip())
      .pipe(gulp.dest(join(cfg.compileDir, cfg.cssDir)))
});



//---------------------------------------------
// Unit Testing
//---------------------------------------------

var testFiles = function() {
    return gulp.src(_.flatten([cfg.vendorFiles.js, cfg.vendorFiles.jsConcat, cfg.testFiles.js, cfg.appFiles.jsunit]));
};

gulp.task('test', ['build-scripts'], function() {
    return testFiles()
      .pipe(karma({
          configFile: cfg.testFiles.config
      }))
});
gulp.task('test-watch', ['build-scripts', 'test'], function() {
    // NOT returned on purpose!
    testFiles()
      .pipe(karma({
          configFile: cfg.testFiles.config,
          action: 'watch'
      }))
});



//---------------------------------------------
// Assets
//---------------------------------------------

// If you want to automate image compression, or font creation,
// this is the place to do it!
//noinspection FunctionWithInconsistentReturnsJS
gulp.task('build-assets-vendor', function() {
    if(cfg.vendorFiles.assets.length) {
        return gulp.src(cfg.vendorFiles.assets, {base: cfg.vendorDir})
          .pipe(gulp.dest(join(cfg.buildDir, cfg.assetsDir, cfg.vendorDir)))
    }
});
gulp.task('build-assets', ['build-assets-vendor'], function() {
    return gulp.src(cfg.appFiles.assets)
      .pipe(gulp.dest(join(cfg.buildDir, cfg.assetsDir)))
});

//noinspection FunctionWithInconsistentReturnsJS
gulp.task('compile-assets-vendor', function() {
    if(cfg.vendorFiles.assets.length) {
        return gulp.src(cfg.vendorFiles.assets, {base: cfg.vendorDir})
          .pipe(gulp.dest(join(cfg.compileDir, cfg.assetsDir, cfg.vendorDir)))
    }
});
gulp.task('compile-assets', ['compile-assets-vendor'], function() {
    return gulp.src(cfg.appFiles.assets)
      .pipe(gulp.dest(join(cfg.compileDir, cfg.assetsDir)))
});



//---------------------------------------------
// Miscellaneous Tasks
//---------------------------------------------

gulp.task('lr-server', function() {
    embedLR = true;
    server.listen(cfg.server.lrPort, function(err) {
        if(err) {
            console.log(err);
        }
        gutil.log('Started LiveReload server');
    });
});

gulp.task('build-clean', function() {
    return gulp.src(cfg.buildDir, {read: false}).pipe(clean());
});

gulp.task('compile-clean', function() {
    return gulp.src(cfg.compileDir, {read: false}).pipe(clean());
});

gulp.task('bump', function() {
    var version = gutil.env['set-version'],
      type = gutil.env.type,
      msg = version ? ('Setting version to '+gutil.colors.yellow(version)) : ('Bumping '+gutil.colors.yellow(type||'patch') + ' version');

    gutil.log('');
    gutil.log(msg);
    gutil.log(gutil.colors.gray('(You can set a specific version using --set-version <version>,'));
    gutil.log(gutil.colors.gray('  or a specific type using --type <type>)'));
    gutil.log('');

    return gulp.src(['./package.json', cfg.bowerJSON])
      .pipe(bump({type:type, version:version}))
      .pipe(gulp.dest('./'));
});