#include <string>
#include <vector>
#include <iostream>
#include <boost/program_options/parsers.hpp>
#include <boost/program_options/variables_map.hpp>
#include <boost/program_options/options_description.hpp>

namespace po = boost::program_options;

int main(int iArgumentCount, char* szArguments[]) {
    // Plus one to avoid putting call of program into arguments
//    std::vector<std::string> arguments(szArguments + 1, iArgumentCount + szArguments);
//    for (const auto& strArgument : arguments) {
//        std::cout << strArgument << std::endl;
//    }
    po::options_description options("Allowed Options");
    options.add_options()
            ("help", "Show Help")
            ("input", po::value<std::vector<std::string>>(), "Set Input Q File");
    po::variables_map variableMap;
    po::store(po::parse_command_line(iArgumentCount, szArguments, options), variableMap);
    po::notify(variableMap);
    if (variableMap.count("input")) {
        auto inputs = variableMap["input"].as<std::vector<std::string>>();
        for (const auto& input : inputs) {
            std::cout << input << std::endl;
        }
    }
}