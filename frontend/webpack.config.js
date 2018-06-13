var HtmlWebPackPlugin = require("html-webpack-plugin");
var CopyWebpackPlugin = require("copy-webpack-plugin");
var HtmlWebpackExternalsPlugin = require("html-webpack-externals-plugin");

module.exports = {
  entry: "./src/index.tsx",
  output: {
    filename: "bundle.js",
    path: __dirname + "/dist"
  },

  // Enable sourcemaps for debugging webpack's output.
  devtool: "source-map",
  mode: "development",

  resolve: {
    // Add '.ts' and '.tsx' as resolvable extensions.
    extensions: [".ts", ".tsx", ".js", ".json"]
  },
  module: {
    rules: [
      // All files with a '.ts' or '.tsx' extension will be handled by 'awesome-typescript-loader'.
      {
        test: /\.tsx?$/,
        use: [
          {
            loader: "babel-loader",
            options: {
              babelrc: false,
              plugins: ["react-hot-loader/babel"]
            }
          },
          "awesome-typescript-loader"
        ]
      },

      {
        test: /\.html$/,
        use: [
          {
            loader: "html-loader"
          }
        ]
      },
      // All output '.js' files will have any sourcemaps re-processed by 'source-map-loader'.
      { enforce: "pre", test: /\.js$/, loader: "source-map-loader" }
    ]
  },

  plugins: [
    new HtmlWebPackPlugin({
      template: "./src/index.html",
      filename: "./index.html"
    }),
    new HtmlWebpackExternalsPlugin({
      externals: [
        {
          module: "react",
          entry: "umd/react.development.js",
          global: "React"
        },
        {
          module: "react-dom",
          entry: "umd/react-dom.development.js",
          global: "ReactDOM"
        },
        {
          module: "semantic-ui-react",
          entry: "dist/umd/semantic-ui-react.min.js",
          global: "semanticUIReact"
        }
      ]
    }),
    new CopyWebpackPlugin([{ from: "assets", to: "assets" }])
  ]
};
