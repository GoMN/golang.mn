/**
 *
 *  Web Starter Kit
 *  Copyright 2014 Google Inc. All rights reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License
 *
 */

'use strict';

// Include Gulp & Tools We'll Use
var config = require('./build.config.js');
var gulp = require('gulp');
var debug = require('gulp-debug');
var glob = require('glob');
var _ = require('lodash');
var $ = require('gulp-load-plugins')();
var rimraf = require('rimraf');
var runSequence = require('run-sequence');
var browserSync = require('browser-sync');
var pagespeed = require('psi');
var reload = browserSync.reload;

// Lint JavaScript
gulp.task('jshint', function () {
    return gulp.src('www/static/scripts/**/*.js')
      .pipe($.jshint())
      .pipe($.jshint.reporter('jshint-stylish'))
      .pipe($.jshint.reporter('fail'))
      .pipe(reload({stream: true, once: true}));
});

// Optimize Images
gulp.task('images', function () {
    return gulp.src('www/static/images/**/*')
      .pipe($.cache($.imagemin({
          progressive: true,
          interlaced: true
      })))
      .pipe(gulp.dest('dist/www/static/images'))
      .pipe(reload({stream: true, once: true}))
      .pipe($.size({title: 'images'}));
});

// Automatically Prefix CSS
gulp.task('styles:css', function () {
    return gulp.src('www/static/styles/**/*.css')
      .pipe($.autoprefixer('last 1 version'))
      .pipe(gulp.dest('dist/www/static/styles'))
      .pipe(reload({stream: true}))
      .pipe($.size({title: 'styles:css'}));
});

// Compile Sass For Style Guide Components (www/static/styles/components)
gulp.task('styles:components', function () {
    return gulp.src('www/static/styles/components/components.scss')
      .pipe($.rubySass({
          style: 'expanded',
          precision: 10,
          loadPath: ['www/static/styles/components']
      }))
      .pipe($.autoprefixer('last 1 version'))
      .pipe(gulp.dest('dist/www/static/styles/components'))
      .pipe($.size({title: 'styles:components'}));
});

// Compile Any Other Sass Files You Added (www/styles)
gulp.task('styles:scss', function () {
    return gulp.src(['www/static/styles/**/*.scss', '!www/static/styles/components/components.scss'])
      .pipe($.rubySass({
          style: 'expanded',
          precision: 10,
          loadPath: ['www/static/styles/']
      }))
      .pipe($.autoprefixer('last 1 version'))
      .pipe(gulp.dest('.tmp/styles'))
      .pipe($.size({title: 'styles:scss'}));
});

// Output Final CSS Styles
gulp.task('styles', ['styles:components', 'styles:scss', 'styles:css']);

// Scan Your HTML For Assets & Optimize Them
gulp.task('html', function () {
    return gulp.src(['www/static/scripts/**/*.html'])
      .pipe(debug({verbose: false}))
        // Minify Any HTML
      .pipe($.minifyHtml())
        // Output Files
      .pipe(gulp.dest('dist/www/static/scripts'))
      .pipe($.size({title: 'client html'}));
});

// Scan Your HTML For Assets & Optimize Them
gulp.task('shtml', function () {
    return gulp.src(['www/templates/*.html'])
      .pipe(debug({verbose: false}))
      .pipe($.useref.assets({searchPath: '{.tmp,www}'}))
        // Concatenate And Minify JavaScript
      //.pipe($.if('*.js', $.uglify()))
        // Concatenate And Minify Styles
      .pipe($.if('*.css', $.csso()))
        // Remove Any Unused CSS
        // Note: If not using the Style Guide, you can delete it from
        // the next line to only include styles your project uses.
      .pipe($.if('*.css', $.uncss({ html: _.flatten([glob.sync('www/templates/*.html'), glob.sync('www/static/**/*.html')]) })))
      .pipe($.useref.restore())
      .pipe($.useref())
        // Update Production Style Guide Paths
      .pipe($.replace('../static/', '/static/'))
        // Minify Any HTML
      .pipe($.minifyHtml())
        // Output Files
      .pipe(gulp.dest('dist/www/templates'))
      .pipe($.size({title: 'server html'}));
});

// Clean Output Directory
gulp.task('clean', function (cb) {
    rimraf('dist', rimraf.bind({}, '.tmp', cb));
});

// Copy application files
gulp.task('app', function () {
    return gulp.src(['www/**/*.go', 'www/*.json', 'www/app.yaml'])
      .pipe(gulp.dest('dist/www'))
      .pipe($.size({title: 'app'}));
});

// Watch Files For Changes & Reload
gulp.task('serve', function () {
    browserSync.init(null, {
        server: {
            baseDir: ['www', '.tmp']
        },
        notify: false
    });

    gulp.watch(['www/static/**/*.html'], reload);
    gulp.watch(['www/static/styles/**/*.{css,scss}'], ['styles']);
    gulp.watch(['.tmp/styles/**/*.css'], reload);
    gulp.watch(['www/static/scripts/**/*.js'], ['jshint']);
    gulp.watch(['www/static/images/**/*'], ['images']);
});

// Build Production Files, the Default Task
gulp.task('default', ['clean'], function (cb) {
    runSequence('styles', 'app', [
        'jshint'
        , 'html'
        , 'images'
        , 'shtml'
    ], cb);
});

// Run PageSpeed Insights
// Update `url` below to the public URL for your site
gulp.task('pagespeed', pagespeed.bind(null, {
    // By default, we use the PageSpeed Insights
    // free (no API key) tier. You can use a Google
    // Developer API key if you have one. See
    // http://goo.gl/RkN0vE for info key: 'YOUR_API_KEY'
    url: 'localhost:8080',
    strategy: 'mobile'
}));
