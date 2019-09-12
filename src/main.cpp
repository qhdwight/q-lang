#include <string>
#include <vector>

#include <boost/program_options/parsers.hpp>
#include <boost/program_options/variables_map.hpp>
#include <boost/program_options/options_description.hpp>

#include <parser/parser.hpp>

namespace po = boost::program_options;

int main(int argc, char* argv[]) {
    po::options_description options("Allowed Options");
    options.add_options()
            ("help", "Show Help")
            ("input", po::value<std::string>()->multitoken(), "Set Input Q File");
    po::variables_map variableMap;
    po::store(po::parse_command_line(argc, argv, options), variableMap);
    po::notify(variableMap);
    ql::parser::Parser().parse(variableMap);
}
