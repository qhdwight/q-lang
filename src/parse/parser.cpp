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

//    std::vector<std::string>

    std::vector<std::string> Parser::extractScopes(std::string code) {
        std::vector<std::string> scopes;
        boost::algorithm::erase_all(code, "\n");
        boost::algorithm::erase_all(code, "\r");
        recurseScopes(code, scopes);
        return scopes;
    }

    void Parser::recurseScopes(const std::string& code, std::vector<std::string>& scopes, int depth) {
        auto level = 0;
        int blockInfoStart = 0, blockStart = 0;
        for (int i = 0; i < static_cast<int>(code.size()); i++) {
            char c = code[i];
            if (c == '{') {
                if (level == 0) {
                    blockStart = i;
                }
                level++;
            } else if (c == '}') {
                level--;
                if (level == 0) {
                    std::string blockWithInfo = code.substr(blockInfoStart, i - blockInfoStart + 1);
                    scopes.push_back(std::move(blockWithInfo));
                    std::string blockContents = code.substr(blockStart + 1, i - blockStart - 1);
                    recurseScopes(blockContents, scopes, depth + 1);
                    blockInfoStart = i + 1;
                }
            }
        }
    }
}
