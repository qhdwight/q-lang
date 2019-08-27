#include "parser.hpp"

#include <iostream>

#include <boost/algorithm/string/erase.hpp>
#include <boost/program_options/variables_map.hpp>

#include <util/read.hpp>

namespace ql::parse {
    std::shared_ptr<ParseNode> Parser::parse(po::variables_map& options) {
        auto sources = options["input"].as<std::vector<std::string>>();
        std::string sourceFileName = sources[0];
        auto src = util::readAllText(sourceFileName);
        auto scopes = extractScopes(src.value());
        for (auto const& scope: scopes) {
            std::cout << scope << std::endl;
        }
    }

    std::vector<std::string> Parser::extractScopes(std::string code) {
        auto level = 0;
        std::vector<std::string> scopes;
        int topOpeningIndex = 0;
        for (int i = 0; i < static_cast<int>(code.size()); i++) {
            char c = code[i];
            if (c == '{') {
                if (level == 0) {
                    topOpeningIndex = i;
                }
                level++;
            } else if (c == '}') {
                level--;
                if (level == 0) {
                    auto contents = code.substr(topOpeningIndex + 1, i - topOpeningIndex - 1);
                    boost::algorithm::erase_all(contents, "\n");
                    boost::algorithm::erase_all(contents, "\r");
                    scopes.push_back(contents);
                }
            }
        }
        return scopes;
    }
}
