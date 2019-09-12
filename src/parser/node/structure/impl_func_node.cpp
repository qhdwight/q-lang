#include <iostream>

#include <boost/algorithm/string/split.hpp>
#include <boost/algorithm/string/trim_all.hpp>
#include <boost/algorithm/string/classification.hpp>

#include <parser/node/implementation/line_node.hpp>

#include "impl_func_node.hpp"

namespace ql::parser {
    void ImplementFunctionNode::parse() {
        // Split by tabs and spaces into tokens, which we use to find what type of node to create
        std::vector<std::string> lines;
        boost::split(lines, m_InnerBody, boost::is_any_of(";"), boost::token_compress_on);
        // Remove first and last blank tokens if they exist
        boost::trim_all_if(lines, [](auto const& token) { return token.empty(); });
        for (auto line: lines) {
            std::cout << line << std::endl;
            std::vector<std::string> tokens;
            boost::split(tokens, line, boost::is_any_of("\t "), boost::token_compress_on);
            boost::trim_all_if(tokens, [](auto const& token) { return token.empty(); });
            m_Children.push_back(std::make_shared<LineNode>(std::move(line), std::move(tokens), weak_from_this()));
        }
    }
}
